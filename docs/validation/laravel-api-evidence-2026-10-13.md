# 項目1〜5の証拠と残件（2026-10-13）

## 提供ログの判定

| 証拠 | 確認できた結果 | 制限 |
| --- | --- | --- |
| [resolveログ](laravel-api-resolve-provided.log) | Composer 2.10.3、PHP 8.3.35/8.4.26/8.5.11。各版の`composer install`はlockからインストールまたは「Nothing to install」、各`check-platform-reqs --no-dev`の全要件がsuccess。プラットフォーム要件を無視するオプションなし | **`resolve.sh` 全体の終了コードは未取得**。ログの最後がsuccessでも0と記載しない |
| `research/laravel-api/composer.lock` と `vendor/composer/installed.json` | 通常パッケージ94件、開発用0件。インストール済み94件と名前・版が全件一致。直接指定と推移依存は[監査](../dependency-audit.md)で区別 | `--no-dev` の範囲のみ。将来の更新版の成立は保証しない |
| [verifyログ](laravel-api-verify-three-versions.log) | 各版で `package:discover`、`migrate`、`metadata:dump`、HTTP検証のコマンド終了コード0。`metadata:dump`は各版1モデル・1123バイトのファイルを生成。各版のHTTPはGET `/api/probe` 200、手動POST 201/GET 200/無効POST 422、API Platform POST 201/返却IRIのGET 200、SQLite行チェックを通過し`PASS` | `verify.py`全体のシェル終了コードはログにない。各版のPASSと最後のstop終了0までは確認。`api-platform:install`は終了0でも後述のファイルキャッシュ書込エラーを出している |
| [監査JSON](laravel-api-vendor-audit.json) と [stderr](laravel-api-vendor-audit.stderr.log) | `audit-vendor.php` 終了コード0はユーザー報告、stderrファイルは0バイト。JSONは94パッケージと89パッケージの構文候補を含み、lockの版と一致 | 構文候補は実行トレースではない。監査コマンドのstdoutへ終了コードは埋め込まれていない |
| [527件の比較](phase2-527-comparison-provided.log) と [nullsafe対象比較](phase2-nullsafe-targeted-provided.log) | 8.3.35/8.4.26/8.5.11でそれぞれ527 matched、0 different、0 errors。対象11件も各版で全件一致 | 収録されたfixtureに限る。PHP全体、Laravel起動を示さない |

提供された `/tmp/laravel-api-resolve.log`、`/tmp/laravel-api-verify-retry2.log`、`/tmp/laravel-api-vendor-audit.json` と空のstderr、`/tmp/nullsafe-targeted.log`、`/tmp/nullsafe-compare-all.log` を、秘密情報らしいトークン・認証情報の有無とサイズを確認して上記パスへ複写した。ログには検証用の固定`APP_KEY`、一時パス、日時が含まれる。resolveログ約692KB・監査JSON約185KB・比較ログ約132KBで、元の判定を再確認できるサイズのため保存した。過去の[依存取得失敗記録](laravel-api-dependency-2026-10-12.md)、[最初のIRI再検証失敗記録](laravel-api-verify-retry-2026-10-13.md)は履歴として残した。後続ログの成功で現在の障害は解消しているが、過去の結果を改変しない。

## キャッシュとIRIに関する事実

`research/laravel-api/verify.py:54-74` は各版を`TemporaryDirectory`で分け、DBを`/probe/database.sqlite`、キャッシュを同じ一時領域の`cache`に置き、`/app/storage/framework/cache/data`へマウントする。`package:discover`がmigration前にメタデータをキャッシュし得るため、`migrate`の後に`shutil.rmtree(cache)`して空ディレクトリを作り直す。`metadata:dump --path=/probe/metadata.dump`の後はファイルの存在と非空を検査し、版ごとのログに1123バイトが残る。`config/api-platform.php:193` の`metadata_dump`はnullなので、この生成ファイルをHTTP処理が読み込んだという証拠ではない。

最初の `/tmp/laravel-api-verify.log` は存在しない`app/ApiResource`への警告、保存先未指定の`metadata:dump`（終了0、生成なし）、API Platform POST時のIRI生成500を記録した。現在の設定は`app_path('Models')`のみ。後続ログではIRI付きPOST 201とGET 200になった。探索先整理とキャッシュ分離・消去のどちらが500解消の決定因子かは、単独介入の比較がないので**未確認**である。後続ログの`api-platform:install`は終了0でも一時キャッシュ内への`file_put_contents`失敗を表示している。HTTP契約は成功したが、このコンソール診断の原因と再現性は残件とする。HTTPログには個別メソッドの実行トレースがない。

## 項目別判定

| 項目 | 今回確定した範囲 | 残件 |
| --- | --- | --- |
| 1 対象と完成条件 | **初期目標を固定**。PHP3版共通のlockとLaravel/API Platform最小アプリで、[完成条件](../laravel-api-target.md)のJSON・SQLite・検証・API作成取得が公式PHPで成立 | Go製PHPでのアプリ起動は未対応。`resolve.sh`全体の終了コードは未確認 |
| 2 実際の依存コード調査 | **固定版・代表実コードの監査を実施**。94件照合、候補棚卸し、起動からHTTP/DB/APIまでの代表コードを[監査表](../dependency-audit.md)に記録 | 全vendorの動的到達・個別拡張APIの使用経路は未確認。網羅的な完了とはしない |
| 3 実行時内部設計 | **今回限定の設計・実装範囲は完了**。[設計](../phase2-design.md)にスコープ、シンボル、ソース文脈、将来の読込方針を記録 | 参照、例外、include/requireの実装は後続 |
| 4 基本文法 | **収録範囲は完了**。null合体、三項、nullsafe、switch/match、ビット演算、文字列展開などの限定範囲が527件の比較で回帰確認 | 未収録の構文・境界やPHP全体は未確認 |
| 5 名前空間と名前解決 | **収録範囲は完了**。namespace/use、定数・関数・クラスの基本名前解決の限定範囲が527件の比較で回帰確認 | 動的名前、複数ファイルautoload、クラス定数などは後続 |

Go実装・fixture・期待値は今回変更しない。次は[依存監査の実装順序](../dependency-audit.md#次のgo実装順序)の最初のファイル読込・autoloadから着手する。Laravel/API Platform全体の動作は今回の完了判定に含めない。
