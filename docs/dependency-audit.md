# Laravel / API Platform 依存コード調査（開始2026-10-12）

現状（2026-10-21）：現行1169件の初回公式PHP3版比較でGenerator委譲、PDO、内部コールバック、SPL境界などに差異を確認し、Go側を修正した。PDO属性の8.5版差も修正した後、公式PHP 8.3.35・8.4.26・8.5.11の版別比較は各1169件一致・差異0・エラー0・終了0となった。Go製PHPの固定アプリは、同じ3設定の実TCPサーバーでLaravel・API Platform・SQLiteの定義済み契約と2000回連続・8並列要求を通過した。[版別生ログと完成判定](full-stack-development.md)。固定vendorの確認済み経路と未収録の依存経路の互換性は区別する。

Generator委譲の境界（修正前の記録）：固定Laravel `LazyCollection`の`yield from`経路を実vendorテストで確認した。`IteratorAggregate::getIterator()`を24段たどる有限の独立fixtureもGoの3設定で実行できるようにした。[実装・暫定期待値](full-stack-development.md)。固定vendorが24段連鎖を実際に使うという証拠ではなく、当時は現行1169件の公式比較は未実施だった。

固定vendorの`Illuminate\Support\LazyCollection`は`flatten()`・`flatMap()`・`chunk()`・`remember()`でGeneratorと`yield from`を使う。Goの[実vendor統合テスト](../cmd/gophp/main_test.go)で、部分消費後の再走査を含む代表経路を3設定で確認した。[ローカル結果・公式比較待ち](full-stack-development.md)。全メソッドや公式PHP3版との一致は未確認。

固定vendorのLaravel Translatorが`array_walk_recursive()`でネストした翻訳文を参照更新する経路をGoの3設定で小さく実行し、`Hi Ada`・`Bye Ada`への置換を確認した。[回帰・制限](full-stack-development.md)。現行1168件の公式比較と、この経路の公式PHP実比較は未実施。

配列走査の格納先（2026-10-21）：固定vendorのSymfony Translation、Laravel Translationなどが`array_walk_recursive()`の参照コールバックを使う。[受け手と添字の副作用fixture](../tests/fixtures/success/array_walk_source_once.php)でGoの配列格納先評価を検査したが、該当vendorの全呼び出し経路や公式PHP3版との一致は未確認。[記録](full-stack-development.md)。

参照取り込みの追加調査（2026-10-21）：固定vendorの`laravel/serializable-closure/src/Serializers/Native.php`は`extract($this->code['use'], EXTR_OVERWRITE | EXTR_REFS)`を使用する。Go側は[小fixture](../tests/fixtures/success/extract_refs_basic.php)で同形のオブジェクト内配列からの参照取り込みを実行したが、固定vendorのシリアライザ自体は実行していない。[対応範囲と公式比較待ち](full-stack-development.md)。

配列走査の追加調査（2026-10-21）：固定vendorのSymfony Translation `XliffLintCommand`は`array_walk()`の参照コールバックを使う。Go製PHPの既存の配列要素参照経路に、名前付き引数を逆順で指定した場合の配列格納先解決を追加し、[独立fixture](../tests/fixtures/success/array_walk_named_reference.php)をGoの3設定でローカル確認した。[詳細](full-stack-development.md)。このSymfonyコマンドの実行や公式PHP3版との一致は確認していない。

固定vendorのPDO束縛（2026-10-21）：Laravel `Connection::bindValues()`は整数に`PARAM_INT`、resourceに`PARAM_LOB`、その他に`PARAM_STR`を指定する。Go製PHPのSQLite束縛へこれらの格納型とstream resourceの基本形を追加し、[対象fixture](../tests/fixtures/success/pdo_bind_value_types.php)と固定アプリ統合テストでローカル確認した。[実装・公式比較待ち](full-stack-development.md)。全vendor経路や公式PHP3版との当時1167件の一致は未確認。

PDOの追加範囲（2026-10-21）：接続とStatementのエラー状態、`errorCode()`／`errorInfo()`、`ATTR_ERRMODE`、列名・取得値に作用する属性の基本経路をGoで確認した。Symfonyの`PdoSessionHandler`にはエラーモード取得、Laravel Connectorには`ATTR_CASE`・`ATTR_ORACLE_NULLS`・`ATTR_STRINGIFY_FETCHES`の設定箇所がある。[状態遷移fixtureと制限](full-stack-development.md)。固定vendorの全PDO操作を実行できたことや公式PHP3版との当時1167件一致は示さない。

SQLite例外の追加経路（2026-10-21）：固定vendorの`SQLiteConnection::isUniqueConstraintError()`が例外メッセージを読み、`Connection::runQueryCallback()`が一意制約例外を選ぶ経路をGoの3設定で実行した。SQLSTATEと`errorInfo`の保持も[対象fixture](full-stack-development.md)で確認した。固定vendorの全DB失敗経路と公式PHP3版での当時1167件一致は未確認。

最新の到達点（2026-10-21）：固定vendorの`Connection::cursor()`とQuery Builder `cursor()`から`LazyCollection`を通る経路をGoの3設定でローカル実行した。生成後にSQLite表を作るfixtureで遅延開始を確認し、SQL失敗時の`PDOException`から`QueryException`への連鎖も捕捉できた。[対象と制限](full-stack-development.md)。これは該当経路のGo実行であり、固定vendorの全DB経路や公式PHP3版での現行1144件一致を示さない。

1167件時点の履歴（2026-10-21）：公式PHP3版の比較済み基準は各997件一致・終了0の[版別生ログ](item11-progress.md)。当時1167件の追加170件と変更後全件は公式比較待ち。[`artisan list`](validation/full-stack-http-api-artisan-list.log)・[`about`](validation/full-stack-http-api-artisan-about.log)・[`route:list`](validation/full-stack-http-api-artisan-route-list.log)、[OpenAPI export](validation/full-stack-http-api-openapi-export.log)、固定アプリのHTTP/SQLite/API Platform基本CRUDとJSON-LD Collection形式の[ページ分割](validation/full-stack-pagination-jsonld.log)はGoのローカル実行で確認した。TermwindのDOM断片処理とSymfony Processの基本経路を通過したが、DOM・process一般の互換を示すものではない。Symfony Responseが`ini_set('default_mimetype', '')`、Requestが`strrchr()`を使う経路、PDOのリクエスト間カーソル寿命まで実行到達した。固定vendorの`Container::resolveFromAttribute()`と`AttributesResourceNameCollectionFactory::create()`は[Goの3設定でローカル実行](validation/full-stack-fixed-vendor-attribute-di.log)した。PHP 8.5側のSymfony定数属性とLaravelの`array_last()`も起動経路で確認した。固定vendorの`ArrayPaginator`が使う`LimitIterator`と`iterator_count()`を確認し、[Goローカルでこの経路を実行](validation/full-stack-array-paginator-local-after.log)した。Symfony Finderの複数探索先が`AppendIterator`を使う経路とSPL子サブクラスの初期値も[Goローカルの対象ケース](validation/full-stack-append-fixtures-local.log)で確認した。Laravel Collectionの`getCachingIterator()`も[Goの固定vendor fixture](validation/full-stack-caching-fixtures-local.log)で実行した。固定vendorのソース上の使用と、Goで通過した実行経路を区別し、全依存経路の対応完了とは扱わない。[詳細](full-stack-development.md)。

固定vendorの`Illuminate\Support\Onceable::tryFromTrace()`は`debug_backtrace(DEBUG_BACKTRACE_PROVIDE_OBJECT, 2)`で呼び出し元のオブジェクトを取得する。従来のGo実装は呼び出し元フレームの`object`を欠き、異なるインスタンスの`once()`キャッシュが混ざった。フレームに実際の受け手を保持する修正後、[固定vendorを使う小fixture](../tests/fixtures/success/laravel_once_object_scope.php)はGoの3設定でオブジェクト別の値を返した。[Eloquentの再帰防止trait](../tests/fixtures/success/eloquent_recursion_cache_vendor.php)も固定vendorを直接読み込み、Goの3設定で再帰値・オブジェクト分離・例外後の解除を確認した。これらは該当経路のローカル実行証拠であり、公式PHP3版との一致や固定アプリ全経路の確認ではない。[詳細](full-stack-development.md)。

固定vendorの`api-platform/state/Util/RequestParser.php`は、クエリ名を16進キーへ変換し、`array_keys()`の結果を型付きクロージャへ`array_map()`で渡す。Goの内部関数コールバックが呼び出し元の`strict_types=1`を誤って適用するため、固定アプリの`?page=1`が500になる経路を[失敗ログ](validation/full-stack-pagination-probe.log)で確認し、共通の弱い型モードへ修正した。[修正後のGo統合ログ](validation/full-stack-pagination-after-callback.log)は31行のページ分割を確認する。公式PHP3版との現行fixture比較は未実施である。

同じ一覧経路のSymfony Serializerは、`is_subclass_of($type, '\\Traversable')`が真の場合にHydraコレクション正規化器を選ぶ。Goの先頭`\\`付きクラス名判定を修正すると、固定vendorの`CachePropertyMetadataFactory::create()`が`serialize(...)`した長文を`hash('xxh3')`へ渡す経路に到達した。240バイト超のハッシュ計算を追加し、[Go統合ログ](validation/full-stack-pagination-jsonld.log)で`Collection`・`totalItems`・`member`の出力を確認した。長文ハッシュは[xxHash公式ベクトル](https://github.com/Cyan4973/xxHash/blob/dev/cli/xsum_sanity_check.c)に対するGo単体テストでも確認した。固定vendorでのこの経路の実行と、公式PHP3版との現在の互換性比較は区別する。

以前の状態（履歴）：型付き参照までの**997件**は[版別生ログ](item11-progress.md)で公式PHP3版と各版一致・終了0。現行**1035件**の追加38件・全件はGoローカル検証済みで、公式比較はPodman接続拒否で未実施。固定vendor内の使用箇所とGoでの実行到達は区別する。実`vendor/autoload.php`はGoで終了0に達したが、`artisan list`は`Illuminate/Container/Container.php:678`の`yield`で停止し、Laravel起動は未達。[全6段階の進捗と実行ログ](full-stack-development.md)。

今回の追加前の比較済み基準：[PHP 8.3](validation/reference-property-974-83-mac.log)・[8.4](validation/reference-property-974-84-mac.log)・[8.5](validation/reference-property-974-85-mac.log)の版別生ログに各974件の個別MATCH、0不一致・0エラー、終了0を確認した。以下の974件時点の「未比較」は当時の履歴。

追加前の比較済み基準：959件は[版別生ログ](item11-progress.md)で公式PHP3版が各959一致・終了0。以下の959件「未比較」は比較前の履歴。

959件比較前（履歴）：呼び出し元スコープの参照も共有判定へ含めた**959件**はGoローカル検証済み、公式PHP3版は環境エラーで未比較。固定vendorでの実行到達は確認していない。[項目11の範囲](item11-progress.md)。

今回の修正前（履歴）：配列参照のコピー境界を修正した**958件**はGoローカル検証済み、公式PHP3版は環境エラーで未比較。固定vendorでの実行到達は確認していない。[項目11の範囲](item11-progress.md)。

958件への修正前（履歴、2026-10-18）：配列要素参照を追加した**955件**はGoローカル検証済み、公式PHP3版はPodman接続拒否で未比較。固定vendorの参照利用経路を実行した証拠とは扱わない。[項目11の対応範囲](item11-progress.md)。

955件追加前の基準（履歴）：変数間参照を加えた**944件**は[版別生ログ](item11-progress.md)で公式PHP3版とも944一致・0不一致・0エラー、終了0。固定vendorの使用確認とGoでの実行到達は引き続き区別し、Composer・Laravel統合の完了とは扱わない。[項目11の対応範囲](item11-progress.md)。以下は944件比較前の履歴。

944件比較前の状況（履歴、2026-10-18）：既存937件は[版別生ログ](items6-10-integration-gate.md)で公式PHP3版と一致、各比較終了0。項目11の変数間参照を加えた944件はGoローカル検証済みで公式比較待ち。固定vendorの使用確認とGo実行到達は引き続き区別する。[項目11の範囲](item11-progress.md)。以下は937件以前の履歴。

937件比較前の検証状況（履歴、2026-10-18）：fixture・期待値は**937件**。`define()`第3引数`null`の診断中断をGo側で修正した。公式PHP3版の対象・全件比較はPodman接続拒否で未実施。固定vendorの使用確認とGo実行到達は区別する。[監査と移行条件](items6-10-integration-gate.md)。以下は936件以前の履歴。

936件時点の検証状況（履歴、2026-10-18）：警告中断経路の追加監査でfixture・期待値は**936件**。`define(null, ...)`の診断ハンドラ例外後に定数が登録されるGo側の経路を修正した。[935件の公式PHP版別生ログ](items6-10-integration-gate.md)は変更前の証拠で、新規1件と修正後全件の公式比較はPodman接続拒否で未実施。固定vendorの使用確認とGoでの実行到達は区別し、composer.lock・vendorは変更していない。[監査と移行条件](items6-10-integration-gate.md)。以下は各時点の履歴。

935件時点の検証状況（履歴、2026-10-18）：fixture・期待値は**935件**（正常515・実行時288・構文133）。公式PHP [8.3.35](validation/items6-10-compound-final-compare-83-mac.log)・[8.4.26](validation/items6-10-compound-final-compare-84-mac.log)・[8.5.11](validation/items6-10-compound-final-compare-85-mac.log)の版別生ログで各版935一致・0不一致・0エラーを確認した。版別比較の終了コード0と`check.sh`のPASSはユーザー報告で、生ログにはない。これは収録fixtureの比較結果であり、固定vendorのGo実行経路への到達やLaravel起動を示さない。公式比較の移行条件は満たしたが、未収録の警告中断経路の監査は残る。[統合検証と項目11への移行判定](items6-10-integration-gate.md)。composer.lockとvendorの版・内容は変更していない。

統合検証追補（履歴、2026-10-17）：固定vendorで確認済みの`call_user_func`、`is_callable`、readonly昇格、enum、相対型の使用例を根拠に、読み込み×型・callableと警告中断を小さいfixtureで組み合わせた。Go側の新規7件はローカル検証済みだったが、vendor実行経路への到達は未確認。この931件の段階ではPodman接続拒否で公式PHP比較待ちであり、[従来の924件版別一致](validation/callable-924-after-fix-summary-mac.log)は当時の変更前の証拠。[後続の比較結果](items6-10-integration-gate.md)を参照。

統合検証前のcallable追補（2026-10-17）：Symfony Console `research/laravel-api/vendor/symfony/console/Attribute/Option.php:54`の`is_callable`と第一級callable、Symfony Routing `Loader/ObjectLoader.php:47`の配列callable検査、Laravel `Illuminate/Session/Middleware/StartSession.php:307`の`call_user_func`、Composer `vendor/composer/autoload_real.php:32,37`の`call_user_func`と`Closure::bind`を固定vendorで確認した。Go側は基本的な文字列／二要素配列／Closureの解決、`is_callable`、`Closure::fromCallable`、第一級callable、対応済み関数に対する`call_user_func`系を部分実装した。`Closure::bind`、`__invoke`、参照、`array_merge`等の追加組み込み関数は残る。これはvendorソースの使用確認であり、Goでその実行経路へ到達した証拠ではない。現在の924件（追加16件を含む）は[ユーザーの版別比較集計](validation/callable-924-after-fix-summary-mac.log)で3版とも一致した。集計に個別fixture行・終了コードはなく、当時の[Podman接続拒否](validation/callable-compare-84-attempt.log)は履歴。[詳細](items6-10-progress.md)。以下の908件の「現在」は当時の履歴として読む。

現在の検証追補（2026-10-16）：相対型と後方型宣言を含む収録908件は[ユーザーの版別集計](validation/special-types-908-after-fix-summary-mac.log)で公式PHP3版それぞれと一致した。追加20件の比較待ちは収録範囲で解消した。集計ログに個別fixture行・終了コードはなく、修正6件と`check.sh`のPASSはユーザー報告である。固定vendorの使用確認とGoによるvendor経路の実行確認は引き続き別であり、後者は未確認。以下の「公式比較待ち」は当時の履歴。[進捗](items6-10-progress.md)。

2026-10-16追補：特殊型の[修正前908件比較](validation/special-types-908-before-fix-mac.log)は、3版とも後方の戻り値型を使う継承2件を拒否し、Goの早過ぎる登録との差を示した。Go側の宣言登録と診断を修正したが、修正後は公式比較待ち。固定vendorの型使用箇所と版の監査事実は変更していない。[進捗](items6-10-progress.md)。

## 相対型・前方参照の追加調査（2026-10-16）

固定vendorの`research/laravel-api/vendor/symfony/error-handler/Exception/FlattenException.php:45`と`vendor/symfony/process/Process.php:206`にnative `: static`、`vendor/symfony/css-selector/Exception/SyntaxErrorException.php:26`にnative `: self`を確認した。`Process.php:286`には`?callable`引数もあり、これは今回の型処理対象外。これらはソース上の使用確認であり、Go製PHPでvendor経路を実行した証拠ではない。Go側は宣言元の`self`／`parent`、called scopeの`static`戻り値、同一ファイルと先行する実行文からの型定義を部分対応とした。新規20件は当時[公式PHP比較試行](validation/special-types-compare-84-attempt.log)がPodman接続拒否で終了2となり、期待値は暫定だった。Composer autoload、一般のcallable／iterable型、vendorの起動経路は未対応・未確認。[進捗](items6-10-progress.md)。

## enumの固定vendor使用確認（2026-10-16）

固定vendorではLaravel `Illuminate/Console/View/TaskResult.php` と `Illuminate/Database/Migrations/MigrationResult.php` がint backed enum、Symfony Console `Output/AnsiColorMode.php`が通常enum、League URI `UriScheme.php`がstring backed enumを宣言する。API Platform JSON Schema `Metadata/Property/Factory/SchemaPropertyMetadataFactory.php`は`BackedEnum`判定、`::cases()`、`->value`を使用する。これはソース中の存在確認であり、Go製PHPでvendor実行経路へ到達した証拠ではない。Goのenum基礎は部分実装である。[修正後の版別比較集計](validation/enum-888-after-fix-summary-mac.log)では3版それぞれ888一致・0不一致・0エラーを確認したが、ログに個別結果とコマンド終了コードはない。[修正前の4件不一致](validation/enum-888-before-fix-mac-summary.log)と[Podman接続拒否](validation/enum-888-diagnostic-compare-all-attempt.log)は履歴である。autoload、Reflection、Laravel実行の残件は維持する。[進捗](items6-10-progress.md)を参照。

## readonlyとプロパティ昇格の固定vendor確認（2026-10-16）

固定lock・vendorの実体では、`laravel/framework` v12.69.3 の `research/laravel-api/vendor/laravel/framework/src/Illuminate/Session/SymfonySessionDecorator.php:18` が `public readonly Session $store` を宣言する。`symfony/type-info` v7.4.20 の `research/laravel-api/vendor/symfony/type-info/Type/TemplateType.php:32-33` と `symfony/serializer` v7.4.20 の `research/laravel-api/vendor/symfony/serializer/Debug/TraceableSerializer.php:35` は、コンストラクタのprivate readonly昇格を使う。これらはvendor内の**存在確認**であり、初期目標の起動・HTTP経路での到達やGo製PHPでの実行成功は未確認。Go側では昇格・readonlyの基本経路を部分実装し、追加27件を含む収録878件は版別比較後の[集計ログ](validation/readonly-878-after-fix-summary-mac.log)で各版一致した。集計ログに個別fixture行・終了コードはない。Reflection、autoload、拡張API等の前提は依然未対応である。[進捗](items6-10-progress.md)を参照。


2026-10-16の動的プロパティ追加直後の記録（履歴）：宣言済み・未宣言のインスタンスプロパティについて、動的な名前式から既存の読み書き・可視性・型検査へ接続した。当時、新規18件はGoローカル検証済みで、公式PHP3版との比較はPodman接続拒否により未実施だった。後続の851件比較では3版各851件の一致を別実行で確認した。[範囲と再比較手順](items6-10-progress.md)。固定vendorの動的プロパティアクセスがGo実行経路で到達するかは未確認であり、Composer autoloadやLaravel起動を対応済みとは扱わない。

2026-10-16の追記：型付きプロパティ／定数の基本実装と複合代入の戻り値修正をGo側へ追加した。[833件の公式PHP3版全件比較](validation/typed-members-833-retry-mac.log)は各版833一致・0不一致・0エラー、終了0で、追加16件の比較待ちは収録fixtureの範囲で解消した。vendorの実行経路到達は未確認。[進捗](items6-10-progress.md)。確認した実コードは`research/laravel-api/vendor/api-platform/hydra/Collection.php:34`の初期値なし`public int $totalItems`と`research/laravel-api/vendor/symfony/service-contracts/ServiceLocatorTrait.php:29`の`private array $loading = []`。型付きクラス定数はvendorの対象検索では見つからなかった。下表の古い状態は監査時点の記録であり、現在の対応判定は本追記と進捗を優先する。

2026-10-15の現在状況：固定vendorの`api-platform/laravel` v5.0.1では`vendor/api-platform/laravel/routes/api.php:34,81,85,94,101`にstaticクロージャのnative `: void`、`Controller/NotExposedController.php:24`にnative `: never`、`Eloquent/Metadata/Factory/Property/EloquentAttributePropertyMetadataFactory.php:60`にprivateメソッドのnative `: never`を目視確認した。Go側に両戻り値型の宣言・制御を部分追加した。817件の[修正後の公式PHP3版比較](validation/void-never-817-after-fix-mac.log)は各版817一致・0不一致・0エラー、終了0。追加30件は収録fixtureの範囲で公式比較済み。ただし、これらのvendor経路がGo製PHPで到達・実行できることは未確認。修正前の[不一致ログ](validation/void-never-817-before-fix-mac.log)は履歴として維持する。

787件時点の履歴：複合型を含む継承メソッドの引数反変・戻り値共変、および宣言時の`int`／`float`判定は、収録fixtureの範囲で[公式PHP3版の成功ログ](validation/type-variance-787-compare-all-mac.log)を確認した。各版787一致・0不一致・0エラー、全件比較終了0。新規18件の比較待ちは解消した。今回は保存済みログの確認で比較を再実行していない。未登録型の後続解決やautoloadを含む個別のvendor実行経路は未確認で、固定vendorの版・監査結果は変わらない。

固定vendorの宣言例：`api-platform/state` v5.0.3の`research/laravel-api/vendor/api-platform/state/ObjectMapper/ObjectMapper.php:21-42`は`ObjectMapperInterface`実装で`object|string|null`引数と`object`戻り値を宣言する。`Provider/ParameterProvider.php:161`には`CollectionOperationInterface&HttpOperation`引数があり、`symfony/uid` v7.4.20の`research/laravel-api/vendor/symfony/uid/Factory/TimeBasedUuidFactory.php:28`には`Uuid&TimeBasedUidInterface`戻り値がある。目視した実コードの宣言事実であり、Go製PHPでの到達・実行成功は確認していない。継承署名の対象範囲は[進捗](items6-10-progress.md)を参照。

765件時点の追補（履歴）：固定`composer.lock`とvendorを確認し、`api-platform/state` **v5.0.3**の`vendor/api-platform/state/Provider/ParameterProvider.php:161`に引数型`CollectionOperationInterface&HttpOperation`、`symfony/uid` **v7.4.20**の`vendor/symfony/uid/Factory/TimeBasedUuidFactory.php:28`に戻り値型`Uuid&TimeBasedUidInterface`を確認した。`symfony/translation` **v7.4.17**の`LoggingTranslator.php:24`にも3構成型のintersectionがあるが、これはreadonly昇格プロパティに付き、今回の引数・戻り値対応の範囲外。`api-platform/serializer` **v5.0.3**の`YamlEncoder.php:25`も同様に範囲外。vendor全体の文字列検索に現れたDNF表記はGuzzleとLaravelのPHPDoc／文字列で、調査した実宣言にnative DNFは確認できなかった。検索だけで使用なしと断定せず、未確認として残す。Go側は引数・戻り値のintersection／DNFを部分実装したが、当時の新規15件は[公式比較待ち](items6-10-progress.md)で、これらのvendor経路のGo実行も未確認。

750件時点の収録fixture（正常405・実行時257・構文88）は、ユーザーのMacでの[修正後全件比較](validation/union-diagnostic-750-after-fix-mac.log)で公式PHP3版とも各750一致・0不一致・0エラー、終了0。union型の構文エラー文言を含む追加23件の比較待ちは収録fixtureの範囲で解消した。これは保存済みログの確認であり、比較の再実行ではない。固定vendorの個別経路での到達は未確認。[現在の進捗](items6-10-progress.md)を参照。

## 718件比較と727件追加直後の記録（履歴）

修正前に比較済みの収録fixtureは718件（正常393・実行時250・構文75）。ユーザーのMacでの[修正後全件比較](validation/item7-dynamic-static-718-after-fix-mac.log)は各版718一致・0不一致・0エラー、終了0。追加11件の公式比較待ちは解消した。今回は保存済みログの確認・保存であり、比較は再実行していない。[修正前の差異](validation/item7-dynamic-static-718-before-fix-mac.log)は履歴である。直接`Class::{throw ...}()`をGoが受理する差は、この718件の比較時点では収録fixture外の既知不具合だった。固定vendorの`Illuminate/Database/Query/Builder.php:39-41`と`Illuminate/Database/DatabaseManager.php:23-25`に複数trait・別名指定、`Illuminate/Database/DatabaseManager.php:119`に`static::`呼び出しがあることを実コードで確認したが、Goでの到達は未確認。過去の監査・比較結果は下に履歴として維持する。

2026-10-15追補（当時）：直接の`Class::{throw ...}()`をGo側の事前検査で拒否し、新規9件を加えて727件となった。[公式PHPとの比較試行](validation/static-throw-compare-all-attempt.log)はPodman接続拒否のため終了2で、この時点では新規ケースの実比較が未実施だった。後続の冒頭727件成功ログで解消した。固定vendorの動的static呼び出しをGoで実際に通せるかは引き続き未確認。

> 2026-10-14のGo側追加状況: `set_error_handler`・`restore_error_handler`・`trigger_error`、エラーレベル定数、既存の一部の警告・非推奨診断の接続は部分実装。追加22 fixtureはGoローカル検証済み。684件のMac比較後に4件を修正し、[再比較ログ](validation/error-handler-684-fix-compare-all-mac.log)で公式PHP3版との全一致を確認した。固定vendorの個別の警告発生経路が実行で到達することは未確認。詳細は[項目6〜10の進捗](items6-10-progress.md)。以下の過去の監査表・固定版は当時の記録として維持する。

662件比較時点で、Go側はローカル読込、継承・interface・trait・staticメンバー・型なしクラス定数、単純な引数・戻り値型、ファイル単位strict_types、値捕捉クロージャ・arrow関数、名前付き引数・配列の文字列キー展開、基本例外制御と例外連鎖を**部分対応**とする。収録fixtureは662件（正常350・実行時237・構文75）。追加21件を含む[修正後の662件成功ログ](validation/throwable-chain-662-compare-all-mac.log)で公式PHP3版とも全一致した。例外連鎖、Throwable取得API、名前付き・可変長引数のトレース、親クラス経由の取得メソッド呼び出しは収録fixtureの範囲で比較済み。今回は保存済みログの確認であり比較は再実行していない。修正前の不一致は[当時のログ](validation/variadic-trace-662-before-fix-mac.log)として維持する。固定vendorの`Illuminate/Foundation/Application.php`が使う宣言・呼び出し全体をGoで解析・実行できることは確認していない。高度な型、参照捕捉・未収録のcallable、Traversable展開、残るThrowable API・エラーハンドラ、Reflection、Composer autoload用組み込み関数と起動経路は残る。[機能別の範囲と未確認事項](items6-10-progress.md)を参照。項目10および項目6〜10全体は部分実装。固定lock・vendorの監査事実と公式PHPでの調査アプリ成功はそのままで、GoによるLaravel起動達成とは区別する。


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
| ローカル`require_once`等のファイル読み込み | Composer生成コード・94件のautoload | `vendor/autoload.php:20`、`vendor/composer/autoload_real.php:29,36,44` | 起動・全HTTP | **部分対応**。ローカル読込・once履歴、通常パスと`include_path`探索・取得・変更・復元は[702件の公式PHP3版比較](validation/include-path-702-after-fix-mac.log)で収録範囲の一致を確認 | Composer autoload全体、呼出し側の型・例外 | P0 | vendor本文を目視。stream wrapper等は対象外、autoload経路のGo実行は未確認 |
| PSR-4、`spl_autoload_register`、`Closure::bind`、`$GLOBALS` | Composer生成コード・94件のautoload | `vendor/composer/autoload_real.php:25-45`、`vendor/composer/ClassLoader.php:381-423` | 起動・全HTTP | **未対応**。上記の`require_once`部分対応だけではComposer生成autoloadは動かない | クロージャ・参照、組み込み関数、ファイルI/O | P0 | 実コードを目視。公式PHPのHTTP成功でautoload入口は確認、個々の内部分岐は未トレース |
| 文字列変数による動的クラス名、二要素配列callable | Laravel 12.69.3、Composer生成コード | `vendor/laravel/framework/src/Illuminate/Collections/Traits/EnumeratesValues.php:818`の`new $class`、`vendor/composer/ClassLoader.php:388`の配列callable登録 | 起動・後続機能 | **部分対応**。`internal/php/parser.go`、`class.go`、`eval.go`で`new $name`・`instanceof $name`・配列メソッド呼出しを処理 | autoload用組み込み関数、クロージャ、全callable形式 | P0 | 収録fixtureは公式PHP3版と比較済み。列挙処理や登録されたcallbackの今回のHTTP経路での到達は未確認 |
| 名前空間、match、null合体 | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:1-40,221-231` | 起動 | 対象部分集合のみ対応。[範囲](phase2-compatibility.md) | 複数ファイルの文脈 | P0 | 本文を目視。未収録の式・型は未確認 |
| 継承・interface・trait | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:36-40`、`vendor/laravel/framework/src/Illuminate/Container/Container.php:25-27` | 起動 | **部分対応**。`internal/php/class.go:56-64,377-521,609-724`に登録・合成・型なし宣言検査があり、収録fixtureは[3版で一致](validation/item7-588-user-comparison-after-fix.log) | 型付き署名の互換性、未収録の競合・循環、autoload | P0 | vendor宣言を目視。Application全体の解析・Go実行は未確認 |
| staticメンバー・静的呼出し・遅延静的束縛 | Laravel 12.69.3、Composer生成コード | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:221-246`、`vendor/composer/autoload_real.php:25-45` | 起動 | **部分対応**。`internal/php/class.go`のセル・呼出しと収録fixtureで固定名の基本形を公式比較済み。動的クラス名・メソッド名と動的staticプロパティを含む追加11件も[718件全一致ログ](validation/item7-dynamic-static-718-after-fix-mac.log)で比較済み。直接`Class::{throw ...}()`は[727件比較](validation/static-throw-727-compare-all-mac.log)で収録範囲の修正を確認済み | 型、autoload用組み込み関数、動的クラス定数など | P0 | `static::`等のコードを目視。vendor起動経路の実到達は未確認 |
| 型なしクラス定数 | Laravel 12.69.3、Symfony http-foundation v7.4.20 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:48`、`vendor/symfony/http-foundation/Request.php:45-65` | 起動・後続 | **部分対応**。`internal/php/class.go:34-44,550-589`の宣言・値と収録fixtureで基本形を確認 | 型付き定数、全定数式、宣言競合 | P0 | vendorの実宣言を目視。各定数が初期経路で参照されるかは未確認 |
| 単純な引数・戻り値型とstrict_types | Laravel 12.69.3 | `vendor/laravel/framework/src/Illuminate/Foundation/Application.php:36-40,221-246` | 起動 | **部分対応**。`internal/php/function.go`、`types.go`、`class.go`、`include.go`で単純型・nullable・ファイル単位の基本検査を実装。収録727件までは3版一致 | 型付きプロパティ／定数、全変換・警告 | P0 | vendor全体の型宣言と版差は未確認 |
| 関数・メソッド引数／戻り値のunion型 | API Platform Laravel v5.0.1、Laravel 12.69.3 | `vendor/api-platform/laravel/Routing/IriConverter.php:100`の`object|string`、`vendor/laravel/framework/src/Illuminate/Filesystem/Filesystem.php:585`の`array|string|int` | APIのIRI生成、ファイル探索 | **部分対応・収録fixtureは公式比較済み**。`internal/php/function.go`、`types.go`、`class.go`にunionの構文・実行時照合と共通の診断表示順を追加。追加23件を含む750件は公式PHP3版で一致 | 完全な型分散、参照引数、vendorの他のPHP構文・組み込みAPI | P1 | 実コードの宣言を目視した。これらのメソッドが初期HTTP経路で実行されたかは未トレース。[公式PHP3版比較](validation/union-diagnostic-750-after-fix-mac.log)はfixtureの範囲を確認 |
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
