# Laravel / API Platform 依存コード調査（開始2026-10-12、現状2026-10-14）

現在、Go側はローカル読込、継承・interface・trait・staticメンバー・型なしクラス定数、単純な引数・戻り値型、ファイル単位strict_types、値捕捉クロージャ・arrow関数、名前付き引数・配列の文字列キー展開、基本的な例外制御を**部分対応**とする。収録fixtureは622件（正常318・実行時230・構文74）。ユーザーがMacで実行した[修正後の公式PHP全件比較](validation/named-exception-622-compare-all-mac.log)は8.3.35・8.4.26・8.5.11の各版622一致・0不一致・0環境エラー、`make compare-all`終了0。605件から追加した17件の比較待ちは解消した。今回は保存済みログの確認・保存であり、比較は再実行していない。[修正前の不一致ログ](validation/named-exception-622-user-before-fix.log)は履歴として維持する。固定vendorの`Illuminate/Foundation/Application.php`が使う宣言・呼び出し全体をGoで解析・実行できることは確認していない。高度な型、参照捕捉・未収録のcallable、Traversable展開、残るThrowable API・エラーハンドラ、Reflection、Composer autoload用組み込み関数と起動経路は残る。[機能別の範囲と未確認事項](items6-10-progress.md)を参照。項目9・10および項目6〜10全体は部分実装。固定lock・vendorの監査事実と公式PHPでの調査アプリ成功はそのままで、GoによるLaravel起動達成とは区別する。


## 過去の実装状況の記録

以下の未実施・未対応は各時点の履歴であり、現在の状況は冒頭の更新を優先する。

547件の公式PHP比較成功を確認した後、Go側は型なしクラスでpublic/protectedメンバーの再宣言・可視性を狭める上書きの拒否、privateメソッドの宣言元への結合、final・abstractクラスとメソッドの基本検査を追加した。追加3 fixtureの公式PHP比較は接続制限で未実施。固定vendorの `Illuminate/Foundation/Application.php:36-40` にある `extends` と `implements`、型宣言、trait等の宣言全体を解析できる状態には至らない。`composer.lock` とvendorは変更していない。現在の機能単位の状態は[進捗](items6-10-progress.md)を参照。

固定vendorの`.php`に対して非対称可視性の`public(set)`等と行頭の`get =>` / `set =>`を再検索した結果は該当0（`rg`終了1）。これは文字列検索であり、PHP 8.4構文の不使用を証明する構文解析ではない。版固有構文はGo側で対応済みと扱わない。

後続レビューで、関数呼び出し・switchを介して親クラスを読み込むケースもGoの単体テストで確認した。これは固定vendorのComposer autoload、interface、trait、型宣言等が動くことを示さない。新規2件の公式PHP比較はPodman接続拒否で未実施である。[項目6〜10の進捗](items6-10-progress.md)を参照。

2026-10-13のレビュー修正では、Go側の単純な `extends` が読込先の親クラスを解決できない不具合と、関数の `static` / `global` 同名宣言の束縛順を修正した。Composer lock・vendorの内容は変更していない。これにより `require` 後の単純継承はGo単体テストで通るが、固定vendorの起動経路に必要なinterface、trait、型、静的呼出し、クロージャ、autoload等は残る。今回の新規fixtureは公式PHPとの実比較前であり、Laravel起動対応や依存全体の互換性として扱わない。[進捗](items6-10-progress.md)を参照。

> 以下の2026-10-12の候補・未実施記述は当時の履歴として残す。現在のlock・vendor実体と3版のHTTP検証に基づく判定は末尾の「2026-10-13 更新」を参照。

選定候補は [Laravel Framework 12.69.3](https://packagist.org/packages/laravel/framework) と [API Platform Laravel 5.0.1](https://packagist.org/packages/api-platform/laravel)。[Laravel のサポート表](https://laravel.com/framework/docs/12.x/releases) は 12 系を PHP 8.2–8.5 対応とし、API Platform の Composer 制約は PHP `>=8.2`、Laravel `^11 || ^12 || ^13` である。PHP 8.3.35 を含む3版に共通する Symfony 7.4 系を明示した。`composer.json` は [調査アプリ](../research/laravel-api/composer.json) に置いた。Composer による解決が環境上未実施のため、**この組み合わせの推移的依存が成立すること、lock の個別パッケージ版、実際の vendor 全体の調査は未確認**。成立しない場合は Composer の衝突原因を保存し、その原因に応じて版別アプリへ分ける。

以下は公開された具体的なファイルを開いてコードを確認した結果である。Laravel は指定タグのファイルを優先した。API Platform のプロバイダも指定タグ v5.0.1 の本文を確認した。ただし vendor にインストールされた実体との照合は未実施である。`起動` はオートロード／コンテナ登録時に必要になる可能性が高いもの、`機能時` は対象の処理に入ったときのものを指す。静的コードの到達判定に依存する箇所は未確認とした。

| パッケージ・ファイルと具体例 | 構文・機能 | Go の現状 | 使用時期・段階 |
| --- | --- | --- | --- |
| [Laravel `composer.json` の `autoload.files`](https://raw.githubusercontent.com/laravel/framework/12.x/composer.json) で `Collections/functions.php` 等を列挙 | Composer autoload、複数ファイル読み込み | 未対応 | 起動・後続 |
| [Laravel `Foundation/Application.php:2-36`](https://raw.githubusercontent.com/laravel/framework/v12.69.3/src/Illuminate/Foundation/Application.php) の `namespace` と `use function ...` | 名前空間、カテゴリ別 import | 今回の部分集合で対応 | 起動・今回。ただしこのファイル全体の解析は不可 |
| 同ファイル `class Application extends Container implements ...` | 継承、interface、型付き引数・戻り値、trait | 未対応 | 起動・後続 |
| 同ファイル `configure()` の `match (true)` と `static::inferBasePath()` | match は今回対応、静的呼び出しは未対応 | 一部 | 起動・今回／後続 |
| 同ファイルの `fn ($path) => ...` と `new static` | arrow function、遅延束縛 | 未対応 | 起動経路次第・後続 |
| [Laravel `Container/Container.php:3-22`](https://raw.githubusercontent.com/laravel/framework/v12.69.3/src/Illuminate/Container/Container.php) の `ArrayAccess`、`Closure`、`ReflectionClass`、`ReflectionAttribute` | 組み込み interface / class、継承、trait | 未対応 | 起動・後続 |
| 同ファイル `getClosure()` の `function (...) use (...)` と `resolve(..., raiseEvents: false)` | クロージャ捕捉、名前付き引数 | 未対応 | コンテナ解決時・後続 |
| 同ファイル `isset($this->bindings[$abstract])`、`??` や配列 | quiet fetch、配列アクセス | 一部対応、今回 `??` 追加 | 起動・今回／後続 |
| [Laravel `Http/Request.php:424-433`](https://raw.githubusercontent.com/laravel/framework/v12.69.3/src/Illuminate/Http/Request.php) の `json_decode($this->getContent() ?: '[]', true)` | 省略三項、`json_decode`、Symfony `InputBag` | 構文のみ今回、関数とクラス未対応 | JSON HTTP 利用時・後続 |
| 同ファイルの `#[\Override]`、`?array` / `: static` | 属性、型宣言、継承 | 未対応 | 起動時にファイル解析が必要・後続 |
| [Laravel `Validation/Validator.php:21-24`](https://raw.githubusercontent.com/laravel/framework/v12.69.3/src/Illuminate/Validation/Validator.php) の `implements ValidatorContract` と `use Concerns\...` | interface、trait | 未対応 | 検証機能時・後続 |
| 同ファイル `validate()` の `throw_if($this->fails(), ...)` | 例外、独自ヘルパー、メソッド呼び出し | 一部。例外とヘルパー未対応 | 検証機能時・後続 |
| [Laravel `SQLiteConnector.php:14-27`](https://raw.githubusercontent.com/laravel/framework/12.x/src/Illuminate/Database/Connectors/SQLiteConnector.php) の `createConnection("sqlite:{$path}", ...)` | 文字列展開は今回対応、PDO 接続未対応 | 一部 | SQLite 利用時・今回／後続 |
| 同ファイル `realpath(...) ?: realpath(base_path(...))` と `$connection->prepare(...)->execute()` | 省略三項、組み込み関数、PDO メソッド連鎖 | 構文の一部のみ今回 | SQLite 利用時・後続 |
| [API Platform `ApiPlatformProvider.php`](https://raw.githubusercontent.com/api-platform/laravel/v5.0.1/ApiPlatformProvider.php) の `namespace`、大量の `use`、`class ... extends ServiceProvider` | 名前解決は今回部分対応、継承なし | 一部。指定タグの本文を確認、vendor 未確認 | 起動・今回／後続 |
| 同プロバイダ `register()` の `static function (Application $app)`、`ReflectionClass`、`class_exists()`、`$app->make()` | 静的クロージャ、型、Reflection、動的サービス | 未対応 | 登録時・後続 |
| 同プロバイダ `boot()` の `Event::listen(...)` と `loadRoutesFrom(...)` | facade 静的呼び出し、ルート読込 | 未対応 | 起動・後続 |

この一覧は vendor の網羅的なスキャン結果ではない。特に API Platform が取り込む `api-platform/{metadata,serializer,state,jsonld,...}`、Symfony 7.4、Composer の生成 autoload の実コードと起動経路は未確認。`composer.lock` 生成後に解決された全パッケージとその実ファイルを列挙し、構文解析による棚卸し、実コードの目視確認、起動時トレースを追加する必要がある。

拡張について、Laravel の必須条件に `ext-ctype/filter/hash/mbstring/openssl/session/tokenizer` が含まれる。[Laravel のパッケージ情報](https://packagist.org/packages/laravel/framework) は PDO を DB 利用時、fileinfo をファイル操作時に挙げている。調査アプリは SQLite 目標に合わせ `ext-pdo_sqlite` と `ext-pdo` を必須にし、`ext-intl` も備える。HTTP JSON は `json_decode` / `json_encode`、検証は `Reflection*` や例外機構を要する。PHP CLI の `php -m` と Composer の実解決ができていないため、その他の推移的拡張要件は未確認である。

## 固定候補の再確認と実コード監査の状態

2026-10-12時点で、[Laravel 12 の公式サポート表](https://laravel.com/framework/docs/12.x/releases) はPHP 8.2–8.5を含む。[Laravel v12.69.3 のソース](https://github.com/laravel/framework/tree/v12.69.3)と[API Platform Laravel v5.0.1 のソース](https://github.com/api-platform/laravel/tree/v5.0.1)の存在を確認した。指定版のパッケージ制約は[Packagist Laravel](https://packagist.org/packages/laravel/framework)、[Packagist API Platform Laravel](https://packagist.org/packages/api-platform/laravel)にあり、API Platform側は PHP `>=8.2`、Laravel `^11 || ^12 || ^13`、Symfony `^7.4 || ^8.0` を許容する。従って直接依存の制約上は3版に共通の候補である。**Composerの全依存解決は未成功なので共通解の確定ではない**。`composer.lock`、`vendor/`、Composerの実バージョン、拡張一覧も未取得である。

次の表は指定タグの公開ソースを目視確認した暫定監査である。パスはインストール後に対応するはずのvendorパスで、行番号はリンク先のタグ付きソースを示す。`起動` / `JSON` / `DB` / `API` は[完成条件](laravel-api-target.md)の段階を表す。到達可能性は実行トレース前には未確認と明記する。

| 機能 | 利用パッケージと固定バージョン | 根拠となるファイルパス・行番号 | 現在のGo実装の対応状況 | 必要になる完成条件 | 実装の前提となる機能 | 優先順位 | 確認方法と未確認事項 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Composer PSR-4 / files autoload、`require` | Composer 2系（実版未取得）、Laravel 12.69.3 | `vendor/autoload.php`（未生成）、[Laravel composer.json](https://github.com/laravel/framework/blob/v12.69.3/composer.json) `autoload` | 未対応 | 起動 | ファイル読込・履歴・名前解決 | P0 | composer.lock生成後にautoload実体を開き、起動トレースで必須範囲を確定。現状未確認 |
| 継承、interface、trait、型付き引数・戻り値、`static::` | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:36-38,221-231` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Foundation/Application.php#L36-L38) | 名前解決のみ一部。継承等は未対応 | 起動 | クラス階層・型・静的ディスパッチ | P0 | クラス宣言とconfigure本文を確認。autoload経由の実行は未確認 |
| `match`、組み込み`is_string`、静的生成 | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:221-231` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Foundation/Application.php#L221-L231) | match部分のみ対応 | 起動 | 静的呼び出し、型宣言、組み込み関数 | P0 | configureの実コードを確認。最初の起動障害の順序は未確認 |
| クロージャ、arrow function、名前付き引数（参照は未確認） | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Container/Container.php:337-340` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Container/Container.php#L337-L340)、`Application.php:243-246` | 未対応 | 起動・API | 変数捕捉、参照セル、引数束縛 | P0 | 本文を目視確認。参照の具体例と各経路の到達は未確認 |
| 例外、ReflectionClass / ReflectionAttribute | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Container/Container.php:16-20,290-314` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Container/Container.php#L290-L314) | 未対応 | 起動・API | try/catch、組み込みReflection、属性 | P0 | 型付きスコープ判定の本文を確認。実際の到達は未確認 |
| `json_decode`、省略三項、Symfony Request | Laravel 12.69.3、Symfony 7.4.x（未解決） | `vendor/laravel/framework/src/Illuminate/Http/Request.php:424-433` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Http/Request.php#L424-L433) | 省略三項のみ対応。JSON関数・HTTP Request未対応 | JSON | HTTP入力、文字列/配列変換 | P1 | 本文を確認。Symfony実版とvendor本文は未確認 |
| PDO SQLite、`realpath`、SQL接続 | Laravel 12.69.3、PHP拡張版は実環境依存 | `vendor/laravel/framework/src/Illuminate/Database/Connectors/SQLiteConnector.php:14-27` [タグソース](https://github.com/laravel/framework/blob/v12.69.3/src/Illuminate/Database/Connectors/SQLiteConnector.php#L14-L27) | 文字列展開の一部のみ。PDO/SQLite未対応 | DB・API | PDO、`ext-pdo_sqlite`、例外、接続設定 | P2 | コネクタ本文を確認。CLIの `php -m` と接続実行は未実施 |
| 属性 `#[ApiResource]`、Eloquent Model | API Platform Laravel 5.0.1、Laravel 12.69.3 | `app/Models/Book.php:1-12`、[API Platform公式例](https://api-platform.com/docs/laravel/) | 属性・継承未対応 | API | 属性解析、Eloquent、Reflection | P3 | 検証アプリに作成。公式PHP・Goでの実行は未実施 |
| サービス登録、静的クロージャ、`class_exists` | API Platform Laravel 5.0.1 | `vendor/api-platform/laravel/ApiPlatformProvider.php:229-248` [タグソース](https://github.com/api-platform/laravel/blob/v5.0.1/ApiPlatformProvider.php#L229-L248) | 未対応 | 起動・API | 継承、クロージャ、動的クラス名、コンテナ | P1 | v5.0.1本文を確認。登録時にどのクロージャが実行されるか未確認 |
| Reflection、ファイル探索、名前解決 | API Platform Laravel 5.0.1 | `vendor/api-platform/laravel/ApiPlatformProvider.php:308-324` [タグソース](https://github.com/api-platform/laravel/blob/v5.0.1/ApiPlatformProvider.php#L308-L324) | 名前解決の一部のみ | API | Reflection、`file_exists`、リソース発見 | P3 | コードを確認。クロージャの到達条件とvendorでの位置は未確認 |
| ルート読込、facade/Event | API Platform Laravel 5.0.1 | `vendor/api-platform/laravel/ApiPlatformProvider.php:1449-1474` [タグソース](https://github.com/api-platform/laravel/blob/v5.0.1/ApiPlatformProvider.php#L1449-L1474) | 未対応 | API | provider boot、イベント、include、HTTPルータ | P3 | 本文を確認。起動時の到達は未確認 |
| 必須拡張 `ctype/filter/hash/mbstring/openssl/session/tokenizer`、調査用`intl/pdo/pdo_sqlite` | Laravel 12.69.3、Composer解決版未確認 | [Laravel composer.json](https://github.com/laravel/framework/blob/v12.69.3/composer.json)、`research/laravel-api/composer.json:6-17` | Go製PHPでは未対応 | 起動・DB・API | 拡張API互換、JSON関数 | P0/P2 | `php -m`、`composer check-platform-reqs`、起動トレースが未実施。個別APIの到達は未確認 |
| Symfony/API Platform推移依存の構文・ランタイム機能 | 7.4.x / 5.0.x（具体版未確定） | `composer.lock`・`vendor/symfony/`・`vendor/api-platform/`（未生成） | 未確認 | 起動・JSON・DB・API | Composer実解決 | P0 | lockの全パッケージ、vendor本文、tokenizer解析、代表ファイルの目視、HTTPトレースが必要 |

依存順は、(1) Composer実解決とvendor監査、(2) `require`/autoload・ファイル履歴、(3) 継承/interface/trait/型・静的アクセス、(4) クロージャ/参照/例外/Reflection、(5) LaravelのJSON HTTPと検証、(6) PDO SQLite、(7) API Platformの属性・リソース・状態処理である。各段階の最小fixtureは順に「同じクラスのautoloadを2回要求」、「親メソッドとinterface実装」、「捕捉変数を持つクロージャと例外処理」、「`GET /api/probe` と422検証」、「SQLiteへ1行保存・取得」、「`#[ApiResource]` の作成・取得」とする。完了条件は各段階のfixtureと対象3版の公式PHPでstdout・stderr・終了コード、HTTP段階ではステータス・本文・DB状態が一致すること。実際のvendor起動経路を調べた後に順序を修正し、フレームワーク全機能は初期目標に含めない。

| 順番 | 作業単位 | 最小fixture | 完了条件 |
| ---: | --- | --- | --- |
| 0 | 共通lockの生成・vendor実体監査 | 3版の `check-platform-reqs` と `audit-vendor.php` | 同一lockを3版でインストールでき、パッケージ版・拡張・起動経路を記録 |
| 1 | ファイル読込・autoload | 2ファイルのPSR-4クラスを2回要求 | クラスを1回だけ定義し、読込履歴・診断が3版と一致 |
| 2 | クラス階層・型・静的アクセス | 親メソッド呼び出しとinterface実装 | 継承・型チェック・静的ディスパッチの正常/失敗例が3版と一致 |
| 3 | クロージャ・参照・例外・Reflection | 捕捉変数を持つクロージャと例外、属性のReflection | コンテナ解決の最小経路が3版と一致 |
| 4 | Laravel HTTP・JSON・入力検証 | `/api/probe` と空titleのPOST | 200 JSONと422 JSONを返し、失敗時のDBに書き込まない |
| 5 | PDO SQLite | 1行の保存・取得と存在しないID | 201/200/404とSQLite行が規定どおり |
| 6 | API Platformの属性・Eloquent・state | `#[ApiResource]` BookのPOST/GET | 201/200のJSON-LDとDB行が規定どおり |

## 2026-10-13 更新：lockと実vendorに基づく判定

直接指定は `research/laravel-api/composer.json` のLaravel 12.69.3、API Platform Laravel 5.0.1、Symfony `type-info:^7.4`・`web-link:^7.4`、PHP `^8.3` と列挙した拡張である。lockの通常パッケージ94件・開発用0件と `vendor/composer/installed.json` の94件は名前・版が全件一致した。`--no-dev` でインストールされた範囲を調べた。固定された推移依存はAPI Platformのmetadata/jsonld/serializer/state等が各v5.0.3、Symfony http-foundation/http-kernel/routing/serializer/validator/type-infoが各v7.4.20、web-linkがv7.4.8、nesbot/carbonが3.14.2。全94件は[監査JSON](validation/laravel-api-vendor-audit.json)にある。

監査JSONは `token_get_all` で89パッケージの構文候補を数えたもの。例えばLaravelのarrow function 834箇所、try 276箇所、API Platform Laravelの属性29箇所は**vendor内の存在数**であり、実行回数ではない。下表は候補から実ファイルを開いた代表例である。3版のHTTP成功は入口から結果までの動作を示すが、関数単位のトレースはなく、個別の内部経路は静的な推定として扱う。パスの `vendor/` は `research/laravel-api/vendor/` を表す。

| 機能 | パッケージと固定版 | 実ファイル・行番号 | 必要な完成条件・段階 | Go側の対応状況と根拠 | 前提となる未対応機能 | 優先 | 確認方法と残る不確実性 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| ローカル`require_once`等のファイル読み込み | Composer生成コード・94件のautoload | `vendor/autoload.php:20`、`vendor/composer/autoload_real.php:29,36,44` | 起動・全HTTP | **部分対応**。`internal/php/include.go:32-109`と[収録比較](validation/item7-588-user-comparison-after-fix.log)はローカル読込・once履歴を確認 | Composer autoload全体、呼出し側の型・例外 | P0 | vendor本文を目視。include_path・stream wrapper等は対象外、autoload経路のGo実行は未確認 |
| PSR-4、`spl_autoload_register`、`Closure::bind`、`$GLOBALS` | Composer生成コード・94件のautoload | `vendor/composer/autoload_real.php:25-45`、`vendor/composer/ClassLoader.php:381-423` | 起動・全HTTP | **未対応**。上記の`require_once`部分対応だけではComposer生成autoloadは動かない | クロージャ・参照、組み込み関数、ファイルI/O | P0 | 実コードを目視。公式PHPのHTTP成功でautoload入口は確認、個々の内部分岐は未トレース |
| 文字列変数による動的クラス名、二要素配列callable | Laravel 12.69.3、Composer生成コード | `vendor/laravel/framework/src/Illuminate/Collections/Traits/EnumeratesValues.php:818`の`new $class`、`vendor/composer/ClassLoader.php:388`の配列callable登録 | 起動・後続機能 | **部分対応**。`internal/php/parser.go`、`class.go`、`eval.go`で`new $name`・`instanceof $name`・配列メソッド呼出しを処理 | autoload用組み込み関数、クロージャ、全callable形式 | P0 | 収録fixtureは公式PHP3版と比較済み。列挙処理や登録されたcallbackの今回のHTTP経路での到達は未確認 |
| 名前空間、match、null合体 | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:1-40,221-231` | 起動 | 対象部分集合のみ対応。[範囲](phase2-compatibility.md) | 複数ファイルの文脈 | P0 | 本文を目視。未収録の式・型は未確認 |
| 継承・interface・trait | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:36-40`、`vendor/laravel/framework/src/Illuminate/Container/Container.php:25-27` | 起動 | **部分対応**。`internal/php/class.go:56-64,377-521,609-724`に登録・合成・型なし宣言検査があり、収録fixtureは[3版で一致](validation/item7-588-user-comparison-after-fix.log) | 型付き署名の互換性、未収録の競合・循環、autoload | P0 | vendor宣言を目視。Application全体の解析・Go実行は未確認 |
| staticメンバー・静的呼出し・遅延静的束縛 | Laravel 12.69.3、Composer生成コード | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:221-246`、`vendor/composer/autoload_real.php:25-45` | 起動 | **部分対応**。`internal/php/class.go:550-600,834-920`のセル・呼出しと収録fixtureで基本形を確認 | 動的クラス・メンバー名、型、autoload用組み込み関数 | P0 | `static::`等のコードを目視。vendor起動経路の実到達は未確認 |
| 型なしクラス定数 | Laravel 12.69.3、Symfony http-foundation v7.4.20 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:48`、`vendor/symfony/http-foundation/Request.php:45-65` | 起動・後続 | **部分対応**。`internal/php/class.go:34-44,550-589`の宣言・値と収録fixtureで基本形を確認 | 型付き定数、全定数式、宣言競合 | P0 | vendorの実宣言を目視。各定数が初期経路で参照されるかは未確認 |
| 単純な引数・戻り値型とstrict_types | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:36-40,221-246` | 起動 | **部分対応**。`internal/php/function.go`、`types.go`、`class.go`、`include.go`で単純型・nullable・ファイル単位の基本検査を実装。収録605件は3版一致 | 複合型、型付きプロパティ／定数、全変換・警告 | P0 | vendor全体の型宣言と版差は未確認 |
| 名前付き引数・引数展開 | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:221-246`、`research/laravel-api/bootstrap/app.php:7-11` | 起動 | **部分対応**。`internal/php/parser.go`、`class.go`で名前付き束縛・整数／文字列キー配列展開を追加 | Traversable展開、束縛エラー時の副作用順序の全境界 | P0 | vendor本文を目視。収録fixtureは公式PHP3版と比較済み、vendor実到達は未トレース |
| クロージャ捕捉・arrow | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Container/Container.php:317-321,403-405` | 起動・ルート | **部分対応**。値捕捉・arrow・staticクロージャの基本形 | 参照捕捉、完全なClosure API | P0 | 本文を目視。参照の実到達は未確認 |
| 例外・Reflection | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Container/Container.php:19-27,999-1001`、`research/laravel-api/bootstrap/app.php:9-10` | 起動・ルート | **例外は部分対応**。`internal/php/exception.go`で基本throw/catch/finallyと一部Error階層を実装。Reflectionは未対応 | 残るThrowable API・エラーハンドラ、Reflection | P0 | 本文を目視。収録fixtureは公式PHP3版と比較済み、vendor経路での到達は未確認 |
| HTTPリクエスト、ルート、レスポンス | Laravel 12.69.3、Symfony http-foundation v7.4.20 | `research/laravel-api/public/index.php:5-9`、`vendor/laravel/framework/src/Illuminate/Foundation/Http/Kernel.php:137,195`、`vendor/laravel/framework/src/Illuminate/Routing/Router.php:749-811`、`vendor/symfony/http-foundation/Request.php:37-42` | GET `/api/probe` 以降 | HTTP実行環境・Request未対応 | 起動、ヘッダー・本文の入出力 | P1 | GET 200を実行確認。内部呼出は未トレース |
| JSON入力・出力、`json_decode`/`json_encode` | Laravel 12.69.3、Symfony http-foundation v7.4.20 | `vendor/laravel/framework/src/Illuminate/Http/Request.php:451-454`、`vendor/laravel/framework/src/Illuminate/Http/JsonResponse.php:76-87`、`vendor/symfony/http-foundation/Request.php:1618-1644` | JSON・検証・API | JSON関数未対応。三項・null合体は部分対応 | HTTP入力、文字列/配列、例外 | P1 | JSON応答200/201/422を確認。特定のJSON関数到達は静的推定 |
| 入力検証と422 | Laravel 12.69.3。Symfony validator v7.4.20はvendor内 | `research/laravel-api/routes/api.php:9-10`、`vendor/laravel/framework/src/Illuminate/Http/Request.php:22`、`vendor/symfony/validator/Mapping/ClassMetadata.php:109-113` | 空bodyのPOST→422 | Laravel検証・例外未対応 | HTTP/JSON、ルール評価、例外変換 | P1 | `errors.title`を実行確認。Symfony validatorのこの経路での使用は未確認 |
| Eloquent/SQLite、PDO、日時 | Laravel 12.69.3、ext-pdo/ext-pdo_sqlite | `research/laravel-api/routes/api.php:11-23`、`vendor/laravel/framework/src/Illuminate/Database/Connectors/SQLiteConnector.php:17-30,40-61` | 保存・取得・API | PDO/SQLite・DB facade未対応 | HTTP、クラス・例外、拡張API | P2 | HTTPとSQLite行を確認。個別PDOメソッドは未トレース |
| 属性探索とReflection | API Platform Laravel v5.0.1、metadata v5.0.3 | `research/laravel-api/app/Models/Book.php:5-11`、`vendor/api-platform/laravel/ApiPlatformProvider.php:321-342`、`vendor/api-platform/laravel/Eloquent/Metadata/Factory/Resource/EloquentResourceCollectionMetadataFactory.php:73-80` | API作成・取得 | **未対応**。クラス継承の部分対応は属性・Reflectionの対応を含まない | 起動、autoload、Eloquent | P3 | BookのPOST/GET成功。Reflectionの実到達は未トレース |
| 状態処理・IRI生成 | API Platform Laravel v5.0.1、state/jsonld各v5.0.3 | `vendor/api-platform/laravel/Eloquent/State/PersistProcessor.php:42`、`Eloquent/State/ItemProvider.php:44`、`Routing/IriConverter.php:100,168-190` | `/api/books` POST→201とIRI GET→200 | 未対応 | 属性、Eloquent、Serializer、URL生成 | P3 | `/api/books/2`の応答を確認。内部関数到達は未トレース |
| メタデータダンプ・キャッシュ | API Platform Laravel v5.0.1 | `vendor/api-platform/laravel/Console/DumpMetadataCommand.php:29,43-56,75-95`、`research/laravel-api/verify.py:59-81` | 調査時の生成 | 未対応 | Reflection、ファイル書込、DBメタデータ | P3 | 3版で1モデル・1123バイト。HTTP時にダンプを読む設定はnullであり別検証 |
| 拡張要件 | 94件のlock、PHP各版 | `research/laravel-api/composer.json:7-22`、`vendor/composer/platform_check.php:5-38` | 起動・DB | PHP拡張一般は未対応 | 各APIの使用箇所ごとの実装 | P0/P2 | 3版の`check-platform-reqs`全行success。ctype/mbstringはpolyfill提供と表示。個別APIの実到達は未確認 |
| GraphQL、Swagger UI、メール、キュー等 | API Platform Laravel v5.0.1、Laravel 12.69.3 | `vendor/api-platform/laravel/ApiPlatformProvider.php:1558-1590`、`vendor/laravel/framework/src/Illuminate/Queue/Worker.php` | 初期目標外 | 未対応 | 対象を決めてから | 対象外 | vendorには存在。今回のHTTP経路では未確認で、必須範囲にしない |

### 次のGo実装順序

拡張APIの具体例として、SQLiteコネクタは `PDO` 接続を作る（`vendor/laravel/framework/src/Illuminate/Database/Connectors/SQLiteConnector.php:17-30`）。LaravelのJSONレスポンスは `json_encode`、Requestは `json_decode` を使う（上表）。vendorの別機能には `mb_strlen`・`openssl_encrypt`・`hash_hmac`（`vendor/laravel/framework/src/Illuminate/Encryption/Encrypter.php:81,108,221`）、`normalizer_normalize`（`vendor/symfony/string/UnicodeString.php:389-409`）、`filter_var`（`vendor/laravel/framework/src/Illuminate/Support/Stringable.php:1519`）もある。これらのファイルの存在は確認したが、今回の6 HTTPリクエストで個々が呼ばれたかは未確認。`ext-ctype`と`ext-mbstring`はComposerの要件表示でpolyfill提供となっており、Go側で拡張全体に対応済みという意味ではない。認証・暗号化や任意のUnicode操作は初期目標外として扱う。

各段階の最小fixtureを同じソースで公式PHP 8.3.35/8.4.26/8.5.11とGo実装へ渡し、CLI段階はstdout・stderr・終了コードをバイト単位で比較する。HTTP段階はステータス・本文・SQLite状態も比較する。現行のGo実装はLaravelを起動できない。

| 順 | 対象範囲 | 最小fixture | 完了条件 |
| ---: | --- | --- | --- |
| 1 | require/require_once、複数ファイル文脈、読込履歴 | 二つのファイルでクラスを定義し二度読む | 定義と診断が3版で一致し、Composerの最小autoloadを読める |
| 2 | 継承、interface/trait、型、静的呼出、名前付き引数 | 親メソッド・interface実装・`Base::make(name: ...)` | 正常/失敗の型・評価順が一致し、bootstrapの構文を解析できる |
| 3 | クロージャ/arrow、捕捉・参照、例外、Reflection/Attribute | 捕捉変数、try/catch、属性のReflection | Composer初期化・Laravelコンテナの最小経路を3版と比較 |
| 4 | 起動に現れた組み込み関数・標準クラス | `spl_autoload_register`と`class_exists`、ファイル操作 | vendor起動の最初の停止点を解消し、次の障害を記録 |
| 5 | HTTP入口・ルート・JSON・検証 | GET `/api/probe`、空title POST | 200/422、本文・DB不変が一致 |
| 6 | PDO SQLite・DB API | 1行の保存/取得と欠落ID | 201/200/404と一時SQLiteが一致 |
| 7 | API Platformの属性・メタデータ・状態処理・IRI | BookのPOSTと返却`@id`のGET | 201/200 JSON-LDとDB2行が一致 |

段階4は1～3の実際の停止点に合わせて小さく進める。認証・キュー・メール・GraphQL等は初期目標に加えない。

## 項目6〜10着手後の補足（2026-10-13）

この2026-10-13追記時点では、Go側にローカルファイルの `include` / `require` / `*_once` と `__FILE__` / `__DIR__`、単純な `extends` と `instanceof`、`global` / `static` ローカル、文字列の変数関数呼び出しを部分実装した。その後、上表のクラス・staticメンバー・型なし定数の部分集合を追加し、収録588件は3版で一致した。Composer生成の `vendor/autoload.php:20-22` は `require_once` の次に静的メソッド `ComposerAutoloaderInit...::getLoader()` を呼ぶ。`vendor/composer/autoload_real.php:5-48` には静的プロパティ・メソッド、`spl_autoload_register`、`Closure::bind`、`$GLOBALS`、`call_user_func`、`require` がある。このため、読込と静的呼出しの部分対応だけではComposer autoloadを起動できない。機能単位の進捗と次の順序は[項目6〜10の進捗](items6-10-progress.md)に記録する。

固定vendorのPHP 8.4固有構文について、`.php` ファイルを対象に `public(set)` 等の非対称可視性と行頭の `get =>` / `set =>` 形式を検索した範囲では該当なし。ただし文字列検索は構文解析による網羅確認ではなく、使用なしとは断定しない。Go側はこれらを8.3設定で無条件に受理せず、現時点では未対応として扱う。調査アプリのキャッシュ書込問題は別件のままであり、今回の言語追加で解消したとは判定しない。
