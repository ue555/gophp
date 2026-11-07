# 新規Laravel 13 + API Platform 5.0.4サンプルをgophp 8.5で起動（2026-10-22、WSL2 Arch Linux）

`~/dev/api-platform-trial`に`composer create-project laravel/laravel:^13.0`（framework v13.38.0）と`composer require api-platform/laravel:^5.0`（v5.0.4、Symfony 8.1系を含む98＋開発33パッケージ）で作った新規アプリ。`app/Models/Book.php`に`#[ApiResource]`と`title`・`author`・`published_year`のmigrationだけを追加した。Composerと`api-platform:install`（設定・Swagger UI資材の公開）は`localhost/gophp-dependency-research:8.5.11`イメージ内の公式PHPで実行し、それ以外はすべてgophpの8.5設定で実行した。固定アプリ（Symfony 7.4系、API Platform 5.0.1）では通らなかった経路が多数あり、以下の修正を本リポジトリの作業ツリーに加えた。gophpenvのv0.4.0エンジンには含まれないため、サンプルの`serve.sh`は`GOPHP_BIN=~/dev/gophp/bin/gophp`で作業ツリーのビルドを指す。

## 到達した障害とgophpの修正

| 障害（発生箇所） | 修正 | 回帰fixture |
| --- | --- | --- |
| `api-platform:install`が`Class "finfo" not found`（Flysystem MIME判定） | 未対応。初回だけ公式PHPで実行 | — |
| `Undefined property: ReflectionProperty::$class`（Symfony Serializer AttributeLoader） | `ReflectionProperty::$class`を追加。`ReflectionClass::getProperties()`が親の宣言を先頭に置き重複して返す不具合、`getProperty()`が親のprivateを返す不具合も修正 | `success/reflection_property_class.php` |
| `Undefined array key 1`（Symfony Finder NumberComparator） | `preg_match()`で後続グループが一致したときの途中の不一致グループを`""`（`PREG_UNMATCHED_AS_NULL`ではnull）で埋める。`PREG_OFFSET_CAPTURE`の不一致は`["", -1]`。名前付きグループは名前→番号の順 | `success/preg_match_unmatched_groups.php` |
| `RecursiveIteratorIterator::setMaxDepth()`未定義（Finder depthフィルタ、Laravel `Filesystem::directories()`） | `setMaxDepth()`／`getMaxDepth()`を追加。最大深さの要素は子へ降りないがhasChildrenは維持（LEAVES_ONLYは出さず、SELF_FIRST／CHILD_FIRSTは出す） | `success/recursive_iterator_max_depth.php` |
| `TypeContextFactory::createFromReflection(): Argument #1 must be of type Reflector, ReflectionMethod given`（Symfony TypeInfo） | `Reflector`インターフェース（Stringable継承）を追加し、ReflectionClass／Enum／FunctionAbstract／Property／Parameter／Attributeに実装。ReflectionNamedTypeはStringableのみ。組み込みReflectionクラスに`__toString`宣言を与え、`ReflectionFunction`を継承するLaravelの`ReflectionClosure`が抽象扱いにならないようにした | `success/reflector_interface.php` |
| `class_implements()`の順序がPHPと異なる（上記fixture作成時に判明） | 親の一覧→宣言順→各宣言インターフェースの一覧を逆順で追加、重複は除外。`autoload=false`時の警告文を`does not exist`に修正 | `success/class_implements_order.php` |
| `ReflectionClass::getNamespaceName()`未定義（TypeInfo） | `getNamespaceName()`・`inNamespace()`、`ReflectionProperty::isPromoted()`を追加 | `success/reflection_namespace_promoted.php` |
| `file()`未定義（TypeInfoが`use`文を収集） | `file()`と`FILE_USE_INCLUDE_PATH`／`FILE_IGNORE_NEW_LINES`／`FILE_SKIP_EMPTY_LINES`を追加。失敗時はPHPと同じ`Failed to open stream: No such file or directory` | `success/file_lines.php` |
| `ReflectionParameter::getDeclaringFunction()`未定義（TypeInfo） | `getDeclaringFunction()`・`getPosition()`を追加 | `success/reflection_parameter_function.php` |

fixtureはいずれもホストの公式PHP 8.5.11と一致し、[全1209件の比較](host-php85-arch-1209-compare-wsl.log)は1209一致・差異0・エラー0、終了0。`go vet`、`go test ./...`（固定アプリの統合テスト含む）、Pythonテスト15件も成功。8.3.35・8.4.26との公式比較はこの環境に該当版がないため未実施で、新規9 fixtureの8.3・8.4期待値は8.5の結果に基づく暫定値（版差のある構文は含まない）。

## サンプルの結果（`GOPHP_BIN=~/dev/gophp/bin/gophp ./serve.sh 127.0.0.1:8090`）

| 要求 | gophp | 公式PHP 8.5.11（`php -S`） |
| --- | --- | --- |
| `artisan key:generate`／`package:discover`／`migrate`（2回）／`cache:clear` | 成功 | — |
| `POST /api/books`（JSON-LD、`publishedYear`含む） | 201、`@id`・`@context`・`publishedYear` | 同じ |
| `GET /api/books/{id}`、`PATCH`（merge-patch） | 200 | 同じ |
| `GET /api/books`（ld+json／json） | Collection、`totalItems` | 同じ |
| `GET /api/docs.jsonopenapi` | OpenAPI 3.2.0、Book schemaに6プロパティ | 正規化後の差分0行 |
| `GET /api/docs`（`Accept: application/ld+json`、Hydra） | 200 | 200。ただしgophp側はErrorリソースのプロパティ（message／code／file等）を欠く（下記） |
| `GET /api/docs`（`Accept: text/html`、Swagger UI） | **500** | 200、42,544B |
| `GET /api/books/999`、`DELETE`→`GET` | 404、204→404 | — |
| 20連続・4並列`GET /api/books` | 全件200 | — |
| title欠落の`POST` | 500（DBのNOT NULL違反） | 同じ500。検証ルール未設定のためAPI Platform本来の挙動 |
| スネークケース`published_year`の入力 | 無視される | 同じ |

## 残る制限（gophp側）

- **Swagger UI**: コンパイル済みBladeの`foreach (...): ... endforeach;`（制御構造の代替構文）をgophpのパーサーが受け付けず`syntax error, unexpected token ":"`。`@foreach`を使うBladeテンプレート全般に影響する。標準welcomeページは使っていなかったため以前の検証では出なかった。
- **Hydraドキュメント**: API Platformの`Error`リソースのプロパティがgophpでは列挙されない。組み込み`Exception`系クラスの`getMessage()`等のアクセサがReflectionから見えないためと推定。OpenAPIには影響しない。
- `fileinfo`（`finfo`クラス）未実装。`api-platform:install`は公式PHPで実行する。
- `ArrayIterator`を継承したユーザークラスが配列を保持しない（`count()`が0、反復が空）。今回の調査で判明、未修正。
- `RecursiveArrayIterator`未実装。
