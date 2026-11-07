# 新規Laravel 12アプリをgophpで起動する試行（2026-10-21）

固定した検証アプリとは別に、`composer create-project laravel/laravel:^12.0`で作った未改変の新規アプリ（laravel/framework v12.70.0）を、WSL2上のArch LinuxでgophpのPHP 8.5設定により起動できるか試した。Composerの実行だけ`localhost/gophp-dependency-research:8.5.11`イメージ内の公式PHPを使い、以降はすべてホストの`gophp`で実行した。配置先は`~/dev/laravel-trial`（Git管理外）。

## 結果

起動は**未達**。到達点と最初に当たった障害は次のとおり。

| 経路 | 到達点 | 障害 |
| --- | --- | --- |
| `vendor/autoload.php`（開発用依存あり） | PHPUnit／Collisionのautoloadファイル | `substr_count()`未実装 |
| `.env`読み込み（vlucas/phpdotenv） | Dotenvの`Regex::split` | `preg_last_error()`未実装 |
| `artisan about`／`migrate`（`--no-dev`、環境変数で設定） | 修正前：`Symfony\Component\EventDispatcher\EventDispatcherInterface`の読み込みで致命エラー。修正後：Console Kernelの起動、コマンド解決まで進み、`console.command`／`console.terminate`イベントの配送で停止 | gophpの不具合（下記）を修正。残りは`krsort()`未実装、例外表示経路の`PCRE_VERSION_MAJOR`未定義 |
| HTTP `GET /up`・`GET /`（`gophp --serve`） | ルーティング、ミドルウェア、ルートクロージャの実行、`Response`生成、Bladeコンパイラの呼び出しまで到達 | `token_get_all()`未実装（tokenizer拡張相当）。500応答 |

migrationは実行されず、SQLiteにテーブルは作成されていない。

## 修正したgophpの不具合

3段以上のインターフェース継承で、中間のインターフェースが祖先のメソッドを互換な形で再宣言している場合（PSR-14 `dispatch(object)` → Symfony Contracts `dispatch(object, ?string = null): object` → Symfony Component）、メソッドを持たない3段目を宣言すると、祖先の宣言を子孫の宣言に対して逆向きに互換性検査し、`Declaration of Psr...::dispatch(object $event) must be compatible with Symfony...::dispatch(object $event, ?string $eventName): object`で致命終了していた。

`interfaceContracts`を、PHPのインターフェースメソッド表の継承規則（親を宣言順に統合し、先に現れた兄弟の宣言を保持、自身の宣言は祖先の宣言を置き換える）に合わせた。既存の`TestInheritedInterfaceContractsAreAllChecked`（ダイヤモンド継承で先の兄弟`Two::load`を保持）と`TestConflictingInterfaceContractsFailAtDeclaration`は維持している。

回帰fixtureは`success/interface_chain_redeclare.php`（再宣言の連鎖と実装クラス、`class_implements`、`ReflectionMethod`の宣言クラス）と`runtime/interface_chain_incompatible.php`（非互換な再宣言の致命エラー、8.5のStack trace付き）。両者はホストの公式PHP 8.5.11と一致し、[全1175件の比較](host-php85-arch-1175-compare-after-interface-fix-wsl.log)も1175一致・差異0・エラー0、終了0。8.3.35・8.4.26との公式比較はこの環境に該当版がないため未実施で、両fixtureの8.3・8.4期待値は8.5の結果と既存fixtureの版差（Stack trace行）に基づく暫定値。

## 未実装の内蔵関数の規模

新規アプリの`vendor`（`--no-dev`）とアプリ本体のPHPソースから関数呼び出し候補を抽出し、公式PHP 8.5.11の内蔵関数一覧と突き合わせた609個のうち、gophpで`function_exists()`が偽になるものは**352個**だった（[一覧](fresh-laravel-12.70-missing-builtins-2026-10-21.txt)）。抽出はメソッド定義名も拾うため上限値であり、起動経路で実際に呼ばれる数ではない。呼び出し箇所が多いものは`mb_substr`・`mb_strlen`・`substr_replace`・`intdiv`・`var_dump`・`substr_count`・`mb_strtolower`・`preg_last_error_msg`・`file`・`array_is_list`・`pathinfo`・`token_get_all`・`random_int`・`array_key_last`など。

## 判断

固定アプリの経路に沿って実装してきたgophpは、未改変の新規Laravelアプリの起動には足りない。固定アプリとの差は主に（1）`.env`読み込み、（2）Symfony EventDispatcher経由のコンソールイベント、（3）Bladeコンパイラ（tokenizer）、（4）文字列・配列・mbstring系の一般的な内蔵関数の網羅性にある。続ける場合は、`preg_last_error`・`krsort`・`substr_count`などの小さな関数から順に実装し、`token_get_all`はgophpの字句解析器を公開する形で設計する必要がある。

## 追記（2026-10-22）：origin/main（v0.4.0）での再試行

リモートの`e1e1ba0`・`602d611`（v0.4.0）を取り込んだ後、同じ`~/dev/laravel-trial`（framework v12.70.0、`--no-dev`）をGoの8.5設定で再試行した。上記の障害はすべて解消している。

| 経路 | 結果 |
| --- | --- |
| `.env`読み込み・`artisan key:generate` | 成功。APP_KEYを`.env`へ書き込み |
| `artisan about` | 成功。Laravel 12.70.0・PHP 8.5.11・Session/Cache/Queue=database を表示 |
| `artisan migrate --force` | 成功。users・sessions・cache・jobs等10テーブルを作成 |
| `gophp --serve --document-root=public` | `/up` 200、`/`（標準welcome、Blade）200で80,650バイト、XSRF／セッションCookie付き、2回目は温キャッシュで同一サイズ。`/nonexistent-page` 404ページ。`/favicon.ico` 200（静的配信）。sessionsテーブルに2行 |

未実装の内蔵関数は`var_dump`・`file`・`mb_substr`・`mb_strlen`・`intdiv`・`array_is_list`などが残るが、未改変アプリの標準経路では呼ばれなかった。この環境のgophpは[gophpenv](https://github.com/ue555/gophpenv)でv0.4.0タグからビルドし、グローバルモード8.5で使用している。
