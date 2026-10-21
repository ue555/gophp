# 項目11：参照セル、変数・配列要素・プロパティの参照

後続段階のPDO状態（2026-10-21）：接続とStatementの独立した`errorCode()`／`errorInfo()`を追加し、Goの3設定で2 fixtureを確認した。現行1161件の変更後全件公式比較は未実施。[段階5の記録](full-stack-development.md)。項目11と全6段階は未完了。

後続段階の追加経路（2026-10-21）：PDO/SQLiteのSQLSTATEとLaravel一意制約例外の基本形をGoで実装し、2 fixtureを加えて当時1144件とした。[段階5の進捗](full-stack-development.md)。項目11のGenerator・反復一般と公式PHP3版の変更後全件比較は引き続き未完了。

最新の連続実装（2026-10-21）：Generatorを通る固定vendorのLaravel/EloquentカーソルをSQLiteで走査し、遅延開始と例外伝播をGoの3設定で確認した。`PDOException`をThrowable階層に追加した。追加3件・当時1142件。公式PHP3版の変更後全件比較は未実施で、期待値は暫定。[実装・ログ](full-stack-development.md)を参照。Generatorの全制御と項目11全体は未完了。

現在（2026-10-21）：現行fixture・期待値は**1161件**（正常737・実行時290・構文134）。参照セル、Generatorの`yield`・`yield from`・遅延実行・基本の`send()`・`throw()`、SPL/ArrayAccessなどをGoで部分実装し、固定アプリのCLI・HTTP経路へ接続した。[継続実装・検証ログ](full-stack-development.md)。公式PHP3版の最後の比較済み基準は各997件一致・終了0で、追加164件と変更後全件の比較はPodman接続拒否により未実施、期待値は暫定。参照・反復一般と項目11全体は未完了。

以前の状態（履歴）：全6段階の連続作業で現行fixture・期待値は**1035件**（正常613・実行時290・構文132）。項目11の参照・反復の入口31件に加え、後続段階のautoload・CLI起動に必要な7件を追加した。公式比較済みの基準は997件のままで、新規38件・全件の3版比較はPodman接続拒否により未実施。項目11では参照引数の格納先拡大、参照戻り値、Generator、ArrayAccess等が残る。[継続実装とログ](full-stack-development.md)。以下の1028件は今回の7件追加前の履歴。

今回追加した`CachingIterator::getCache()`の値コピーと配列構文による更新・削除の[回帰fixture](../tests/fixtures/success/caching_iterator_array_copy.php)はGoの3設定で通過した。固定アプリのAPI Platform CRUDも、既存の8.4設定に加え[8.3・8.5設定の統合テスト](validation/full-stack-api-version-matrix-local.log)で通過した。公式PHP3版の現行1134件との比較は未実施で、項目11の全反復境界は未完了。

## 2026-10-21 のCachingIterator境界修正

先読み時に内側の`current()`を`key()`より先に呼ぶ順序へ修正し、コンストラクタと`setFlags()`で内部状態ビットを公開フラグとして受け取らないようにした。[3件の新規fixtureと結果](full-stack-development.md)、[ローカル対象ログ](validation/full-stack-caching-boundary-local.log)に記録した。総1132件で、公式PHP3版の新規3件と変更後全件は未比較。`CachingIterator`全APIと反復中変更は引き続き未完了。

## 2026-10-21 のCachingIterator

固定vendorのLaravel Collectionが生成する`CachingIterator`の基本形を追加した。現在値と内部Iteratorの先読み位置を別々に保持し、`hasNext()`と`FULL_CACHE`を接続する。追加5件・総1129件（正常705・実行時290・構文134）は[Go4設定の対象結果](validation/full-stack-caching-fixtures-local.log)、[全件test](validation/full-stack-caching-test.log)・[vet](validation/full-stack-caching-vet.log)・[build](validation/full-stack-caching-build.log)・[差分](validation/full-stack-caching-diff.log)で検証した。公式PHP3版の変更後全件と新規対象は未比較で、各版0件実比較、1129件未実施。期待値は暫定である。全`setFlags()`診断や`RecursiveCachingIterator`、参照を伴う反復は残件とする。[全6段階の進捗](full-stack-development.md)にMac再比較手順と実アプリ経路を記録する。

## 2026-10-21 のAppendIteratorとSPL子インスタンス

`AppendIterator`を内部Iteratorの列として表し、共通の`iteratorMethod`から順に委譲する。Finderが複数の`in()`先を指定する経路で、SPLの`getChildren()`がサブクラスの型付きプロパティ初期値を消していたため、具体クラスの生成とコンストラクタ呼び出しへ変更した。追加5件・当時総1124件（正常700・実行時290・構文134）の[Go3設定の対象結果](validation/full-stack-append-fixtures-local.log)は15/15件成功。[全件test](validation/full-stack-append-test.log)・[vet](validation/full-stack-append-vet.log)・[build](validation/full-stack-append-build.log)・[差分検査](validation/full-stack-append-diff.log)と[`artisan list`](validation/full-stack-append-artisan-list.log)も終了0。公式PHP3版との対象・全件比較はPodman接続拒否が継続し未実施で、各版0件比較、1124件未実施。期待値は暫定。反復中のAppendIterator変更とSPL子クラスの全境界、参照・反復一般は残る。Macの正確な版別比較コマンドは[全6段階進捗](full-stack-development.md)に記載した。

## 2026-10-21 のIteratorラッパー

`IteratorIterator`・`LimitIterator`の基本形と`iterator_count()`を追加した。固定vendorの`ArrayPaginator`を実行し、ページの1件列挙と件数をGoローカルで確認した。[実装範囲・6 fixture・ログ・未検証境界](full-stack-development.md)。公式PHP3版ではまだ比較していないため、項目11全体は未完了。

## 2026-10-21 の通常オブジェクト反復

通常オブジェクトの可視プロパティに対する値渡し・参照渡しの`foreach`基本形を追加した。privateスロット、protectedスコープ、型付き未初期化と参照セルの制約を既存格納先へ接続した。[実装・4 fixture・検証・Mac再比較手順](full-stack-development.md)を参照。反復中の変更などは未検証で、公式PHP3版の実比較待ちを維持する。

## 2026-10-21 の参照Generator宣言検査

参照Generatorに`yield from`がある場合、到達したときのGo独自エラーから、未呼び出しでも事前に拒否するFatal診断へ変更した。入れ子の通常クロージャは独立して検査する。[PHP公式バグ記録](https://bugs.php.net/bug.php?id=71252)を根拠とするが、対象3版の実比較は未実施。[Go3設定の対象ログ](validation/full-stack-reference-yield-from-local.log)と[詳細](full-stack-development.md)。現行1106件のうち追加109件の期待値は暫定。項目11は部分実装。

## 2026-10-21 の通常関数・メソッド・Closure参照戻り値（1104件時点）

通常関数、インスタンス／staticメソッド、Closure、アロー関数の参照戻り値を、`=&`による束縛と通常代入による値コピーに分けた。元の変数・配列要素・プロパティ・staticローカルと共有する7 fixtureを追加し、当時1104件（正常681・実行時290・構文133）。[Goの3設定での出力・終了コード](validation/full-stack-reference-return-local.log)と[`artisan list`](validation/full-stack-reference-return-artisan-list.log)はローカル検証であり、公式PHP3版の一致証拠ではない。参照を返す呼び出しを参照引数へ一度だけ評価して渡す基本形も追加した。非変数式や括弧による戻り値の境界は残る。[設計と再比較手順](full-stack-development.md)。最後の公式比較済み997件から追加107件の期待値は暫定で、項目11は部分実装のままである。

## 2026-10-21 の参照Generator（参照戻り値追加前の履歴）

`function &generator()`の`yield`が変数・配列要素・通常プロパティを参照セルとして保持し、Generatorまたは`IteratorAggregate`からの参照foreachがそのセルへ束縛する基本形を追加した。通常の`current()`や値渡しforeachでは値だけを渡す。[追加5件のGo3設定での実行](validation/full-stack-generator-reference-local.log)と[test](validation/full-stack-generator-reference-test.log)等のローカル検証は成功。公式比較は[Podman終了125](validation/full-stack-generator-reference-podman.log)のため未実施で、期待値は暫定。参照`yield from`、式の参照yield、通常関数の参照戻り値、一般の反復中変更は残る。[詳細とMac比較手順](full-stack-development.md)。

## 参照捕捉・参照引数・参照foreach・反復interface（2026-10-19、1028件時点の履歴）

[公式PHP 8.3](validation/reference-typed-warning-997-83-mac.log)・[8.4](validation/reference-typed-warning-997-84-mac.log)・[8.5](validation/reference-typed-warning-997-85-mac.log)の版別生ログはそれぞれ997件の個別MATCH、997一致・0不一致・0エラー、`exit_code: 0`を記録する。型付き参照の警告例を実行時失敗系へ移した後の結果である。下の「997件比較待ち」は当時の履歴として残す。

現行fixture・期待値は**1028件**（正常606・実行時290・構文132）。比較済み997件は維持し、項目11と段階2・3の入口機能を含む新規31件を追加した。Go側ではクロージャの`use (&$x)`、既存変数を渡す関数・メソッド・クロージャの`&$parameter`、配列変数の参照foreach、Iterator／IteratorAggregate／Traversableの基本反復、Countableの基本呼び出しを実装し、ローカルの期待値テストで検証した。参照セルは外側の変数・クロージャ・引数の束縛を共有し、参照foreachでは各要素のセルを順に値変数へ結ぶ。IteratorとCountableは既存のメソッド呼び出しへ接続した。参照と反復の統合fixtureもGoで`2|5`を出した。[全6段階の依存順と制限](full-stack-development.md)を参照。

Podmanソケットは接続拒否（`podman ps`終了125）のため、追加31件と全1028件の公式PHP3版比較は**未実施**であり、期待値は暫定。配列要素・プロパティ・未定義変数・展開・可変長の参照引数、参照戻り値、一般オブジェクトのpublicプロパティ反復、ArrayAccess、Generator、反復中の全ての変更境界は未完了。項目11と全6段階を完了扱いにしない。Macでの版別の対象・全件比較コマンドは[全体進捗](full-stack-development.md)に記録した。

## 2026-10-20 の全6段階継続作業

参照・反復の既存部分を使い、Generatorの基本遅延実行から固定vendorのautoload、Laravel CLI、HTTP/SQLite、API Platformの基本CRUDまでローカルの実行経路を進めた。当時のfixture・期待値は1092件、Goの`make test`・`make vet`・`make build`・`git diff --check`は終了0。[統合実行と残件](full-stack-development.md)を参照。項目11の参照・反復一般は部分実装のまま。最後に公式比較したのは997件で、追加95件と変更後全件はPodman接続拒否により3版とも未比較である。`Generator::send()`・`throw()`の基本形に加え、Generator間の委譲転送とIterator／IteratorAggregateの基本委譲をGoローカルで検証した。参照を含む反復変更、Traversable境界と長時間継続は未対応・未検証として残す。新規4件と変更後全件の公式比較は未実施である。

## 型付きプロパティへの参照（2026-10-19、997件比較前の履歴）

fixture・期待値は**997件**（正常576・実行時289・構文132）。比較済み974件のソース・期待値は維持し、型付き通常／staticプロパティ、別名からの通常代入・複合代入・増減・`??=`、配列要素との共有、複数制約、再束縛、`unset`、未初期化nullable、readonly拒否を23件追加した。参照セルは型付きプロパティの束縛を制約源として保持する。変数・配列要素・通常／staticプロパティからの書き込みは共通の検査を通り、書き込み箇所の`strict_types`を適用する。複数の制約が異なる値への暗黙変換を要求する場合はTypeErrorで拒否する。プロパティの付け替え・`unset`では旧セルの制約を外す。クロージャが保持する`boundThis`も到達可能性に含め、`$box`変数を`unset`してもクロージャ経由で生存するオブジェクトの型制約と、配列コピー時の参照共有を維持する。この境界を2件のfixtureに追加した。循環参照とGCの細かな時点は未検証。

[8.3](validation/reference-boundthis-997-83-mac.log)・[8.4](validation/reference-boundthis-997-84-mac.log)・[8.5](validation/reference-boundthis-997-85-mac.log)の前回の版別生ログでは、各996一致・0不一致・1エラー、比較終了2。`success/reference_typed_warning.php`は公式PHPが終了255であり、成功fixtureという分類が誤っていた。公式stderrに基づいて同じソースを`runtime/`へ移し、8.3のクロージャ名表示差を版別期待値に記録した。Goでは、参照への損失を伴う型変換の診断ハンドラが例外を投げた場合、その例外を先行原因として参照型のTypeErrorを発生させる。ハンドラ例外を保留したまま返す経路も修正した。前回ログは移動前のパスに対する証拠であり、移動後のfixtureと修正後の全件一致を証明しない。

修正後の[対象8.3](validation/reference-typed-warning-fix-target-83.log)・[8.4](validation/reference-typed-warning-fix-target-84.log)・[8.5](validation/reference-typed-warning-fix-target-85.log)、[全件8.3](validation/reference-typed-warning-fix-all-83.log)・[8.4](validation/reference-typed-warning-fix-all-84.log)・[8.5](validation/reference-typed-warning-fix-all-85.log)はPodman接続拒否（終了125）により公式PHP起動前に比較器が終了2。各版0一致・0不一致、対象1件と全997件の公式比較は未実施。[test](validation/reference-typed-warning-fix-test.log)・[vet](validation/reference-typed-warning-fix-vet.log)・[build](validation/reference-typed-warning-fix-build.log)・[差分](validation/reference-typed-warning-fix-diff.log)は終了0。移動後の期待値は公式出力を新しいパスへ適用した暫定値として、公式再比較を待つ。

前回のGo側検証履歴：[test](validation/item11-typed-reference-bound-this-test.log)・[vet](validation/item11-typed-reference-bound-this-vet.log)・[build](validation/item11-typed-reference-bound-this-build.log)・[差分](validation/item11-typed-reference-bound-this-diff.log)は各終了0。[Goの3設定](validation/item11-typed-reference-bound-this-local-versions.log)で当時の新規23件は各23件が暫定期待値と一致。Podmanは[接続拒否](validation/item11-typed-reference-bound-this-podman.log)で終了125。[対象8.3](validation/item11-typed-reference-bound-this-target-83.log)・[8.4](validation/item11-typed-reference-bound-this-target-84.log)・[8.5](validation/item11-typed-reference-bound-this-target-85.log)と[全件8.3](validation/item11-typed-reference-bound-this-all-83.log)・[8.4](validation/item11-typed-reference-bound-this-all-84.log)・[8.5](validation/item11-typed-reference-bound-this-all-85.log)は公式PHPの起動前に比較器が終了2。この時点の暫定期待値の一部を上記の公式PHP生ログに基づいて訂正した。項目11全体は部分実装。

Macでの対象比較は下記をリポジトリルートで版ごとに実行する。`83 8.3.35 8.3`を`84 8.4.26 8.4`、`85 8.5.11 8.5`に替える。各fixtureを比較し、差異を修正した後に全件を比較する。`make compare`はローカルtest・buildも実行する。各ログには実コマンドと直後の終了コードを残す。

```sh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r suffix version target <<< "$spec"
  for fixture in tests/fixtures/success/reference_typed_*.php; do
    log="docs/validation/item11-typed-reference-$(basename "$fixture" .php)-${suffix}-mac.log"
    printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures %s --interpreter ./bin/gophp --php-version=%s {file}\ntarget_php: %s\n' "$suffix" "$version" "$fixture" "$target" "$version" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${suffix}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
  done
  log="docs/validation/item11-typed-reference-all-${suffix}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\ntarget_php: %s\n' "$suffix" "$version" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${suffix}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
done
```

次の実装単位は参照引数・参照戻り値、クロージャ参照捕捉、参照foreachとTraversable。参照による配列循環、clone、マジックメソッド、GC時点の完全再現も今回の範囲外。型付きプロパティ参照の初期化・複数制約の診断と版差は公式3版比較で確定させる。

## 974件の公式比較済み基準（今回の追加前）

[8.3.35](validation/reference-property-974-83-mac.log)・[8.4.26](validation/reference-property-974-84-mac.log)・[8.5.11](validation/reference-property-974-85-mac.log)の版別生ログは、それぞれ974件の個別MATCH、**974一致・0不一致・0エラー**、`exit_code: 0`を記録する。959件から追加した型なしプロパティ参照15件は収録fixtureの範囲で公式比較済み。版ごとの`make compare`であり、`make compare-all`成功記録ではない。`check.sh`のPASSはユーザー報告。以下の974件時点の「未比較」は当時の履歴。

## 型なしプロパティへの参照（2026-10-18、974件比較前の履歴）

fixture・期待値は**974件**（正常554・実行時288・構文132）。既存959件は変更せず、通常／staticプロパティ、変数・配列要素との共有、参照の付け替え、`unset`後の寿命、動的名の評価、private／protectedの可視性、static継承、警告ハンドラ例外、配列コピー時にプロパティ内へ残る参照を15件のfixtureで確認した。`referencePlace`が受け手・名前・classを一度だけ解決し、既存の`propertyKey`と`staticPropertyLocation`で可視性を検査する。通常書き込みは既存セルの値を更新し、`=&`だけが対象プロパティの束縛を別セルへ付け替える。通常の読み取りやstaticアクセスは内部セルを値として露出させない。型付きプロパティの参照一般は未実装で、Go固有の明示的なエラーで拒否する。readonlyは参照による間接変更を拒否する。型付き参照セルの制約登録・共有解除時の制約管理を次の実装単位へ残す。

[test](validation/item11-property-reference-test.log)・[vet](validation/item11-property-reference-vet.log)・[build](validation/item11-property-reference-build.log)・[差分](validation/item11-property-reference-diff.log)は終了0。[Goの3設定](validation/item11-property-reference-local-versions.log)で新規15件は各15一致。Podmanは[接続拒否](validation/item11-property-reference-podman.log)で終了125。[対象8.3](validation/item11-property-reference-targeted-83.log)・[8.4](validation/item11-property-reference-targeted-84.log)・[8.5](validation/item11-property-reference-targeted-85.log)と[全件8.3](validation/item11-property-reference-all-83.log)・[8.4](validation/item11-property-reference-all-84.log)・[8.5](validation/item11-property-reference-all-85.log)はいずれも公式PHPの起動前に比較器が終了2。各版0一致・0不一致・1環境エラー、974件未実施。新規15件の期待値と診断は**暫定**で、実装単位の公式比較完了判定は保留する。項目11全体も部分実装。

Macでの再比較はリポジトリルートで、次を版ごとに順番に実行する。対象15件を個別比較し、差異を修正してから全件を比較する。`make compare`はローカルtest・buildも実行する。各ログにはコマンドと直後に得た終了コードを追記する。

```sh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r suffix version target <<< "$spec"
  for fixture in tests/fixtures/success/reference_property_*.php; do
    log="docs/validation/item11-property-$(basename "$fixture" .php)-${suffix}-mac.log"
    printf 'command: compare_php.py --php ./containers/php-version/bin/php%s --php-version %s --fixtures %s --interpreter ./bin/gophp --php-version=%s {file}\ntarget_php: %s\n' "$suffix" "$version" "$fixture" "$target" "$version" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${suffix}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
  done
  log="docs/validation/item11-property-all-${suffix}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\ntarget_php: %s\n' "$suffix" "$version" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${suffix}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
done
```

## 959件の公式比較（2026-10-18、今回の追加前の基準）

[8.3](validation/reference-scope-959-83-mac.log)・[8.4](validation/reference-scope-959-84-mac.log)・[8.5](validation/reference-scope-959-85-mac.log)の版別生ログに各959件の個別MATCH、**959一致・0不一致・0エラー**、`exit_code: 0`を確認した。944件から追加した配列要素参照関連15件と訂正した期待値は、収録fixtureの範囲で公式比較済み。版ごとの`make compare`であり、`make compare-all`の成功記録ではない。`check.sh`のPASSはユーザー報告。以下の959件時点の「公式比較待ち」は当時の履歴である。

## 呼び出し元の参照を共有判定へ追加（2026-10-18、959件比較前の履歴）

fixture・期待値は**959件**（正常539・実行時288・構文132）。新規`reference_array_caller_scope.php`は、呼び出し元のローカル変数が配列要素への参照を保持したまま、その配列を関数へ値渡しする。従来のGo実装は呼び出し先のスコープだけを探索して要素セルを分離し、`1|1`を出した。呼び出し時に退避したスコープを実行中のスタックに保持し、配列書き込み時の共有判定に含めた結果、3設定とも`9|9`になる。関数から戻る際に退避スコープの登録を外す。既存958件のfixture・期待値は変更していない。[PHP公式マニュアル](https://www.php.net/manual/en/language.references.whatdo.php)は、配列を値渡しした関数呼び出しでも要素内の参照を保持すると説明するが、**この再現例の公式PHP実比較は未実施**であり、新規期待値は暫定である。

`make test`（Python12件を含む）・`make vet`・`make build`・`git diff --check`は終了0（[test](validation/item11-array-caller-scope-test.log)・[vet](validation/item11-array-caller-scope-vet.log)・[build](validation/item11-array-caller-scope-build.log)・[差分](validation/item11-array-caller-scope-diff.log)）。[対象8.3](validation/item11-array-caller-scope-targeted-83.log)・[8.4](validation/item11-array-caller-scope-targeted-84.log)・[8.5](validation/item11-array-caller-scope-targeted-85.log)、[全件8.3](validation/item11-array-caller-scope-all-83.log)・[8.4](validation/item11-array-caller-scope-all-84.log)・[8.5](validation/item11-array-caller-scope-all-85.log)は、Podman接続拒否（公式PHP起動終了125）により比較器・`make compare`が終了2。各版0一致・0不一致・1環境エラー・959件未実施。項目11全体は部分実装のままで、今回の修正単位の公式比較完了判定は保留する。Macで再実行する際は、下記の版別コマンドの対象を`tests/fixtures/success/reference_array_caller_scope.php`にし、全件は`make compare`を版ごとに実行する。

## 配列コピー時の参照分離を修正（2026-10-18、今回の修正前の履歴）

fixture・期待値は**958件**（正常538・実行時288・構文132）。新規3件は`unset($x)`後に配列をコピーし、要素への代入、追記後の代入、コピー元への代入を分けて確認する。既存`reference_array_copy_after.php`は、コピー後でも**書き込み前**に外部別名をunsetするため、暫定期待値を`4|4|6`から`1|4|4`へ訂正した。以前の期待値はGoの出力に依拠しており、今回の修正根拠は[PHP 8.4の`zend_array_dup_value()`](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_hash.c)である。そこでは参照セルの参照数が1で、参照先が複製元配列自身ではない場合、セルの値をコピーする。**公式PHPの実出力はこの環境では未確認**なので、訂正した期待値と新規3件の期待値はいずれも暫定とする。既存fixtureのソースは変更していない。

Go側は配列を書き込む際、同じセルを保持する別の変数・配列要素・static格納先が残るか調べる。同じ配列オブジェクトを複数の変数が指していても、その中の1要素を重複して数えない。共有先のないセルは新しい配列では通常値に戻し、共有先のあるセルは維持する。これを要素代入だけでなく、追記・要素の付け替え・unsetに通した。したがって「コピー前に参照があれば常に共有」という前回の説明は撤回する。クロージャ以外の退避済み呼び出しフレームや未対応の間接循環まで含む完全なPHP参照カウントではないため、その境界は後続で確認する。

`make test`（Python12件を含む）、`make vet`、`make build`、`git diff --check`の結果は[今回の検証ログ](validation/item11-array-copy-fix-test.log)以下に記録する。公式PHP3版の[対象8.3](validation/item11-array-copy-fix-targeted-83.log)・[8.4](validation/item11-array-copy-fix-targeted-84.log)・[8.5](validation/item11-array-copy-fix-targeted-85.log)と[全件8.3](validation/item11-array-copy-fix-all-83.log)・[8.4](validation/item11-array-copy-fix-all-84.log)・[8.5](validation/item11-array-copy-fix-all-85.log)はPodman接続拒否で終了2。各版0一致・0不一致・1環境エラー・958件未実施であり、修正単位の完了判定は公式3版比較後に行う。Macでの正確な版別対象・全件比較コマンドは次節の手順を利用する。今回の新規fixtureも`reference_array_*.php`に含まれる。

## 配列要素への参照（2026-10-18、今回の修正前の履歴）

fixture・期待値は**955件**（正常535・実行時288・構文132）。944件の公式PHP3版一致は下記の比較済み基準であり、今回追加した`reference_array_*.php`11件と修正後の全955件は**公式比較待ち**。新規期待値は暫定。`make test`（Python 12件を含む）、`make vet`、`make build`、`git diff --check`は各終了0（[test](validation/item11-array-test.log)・[vet](validation/item11-array-vet.log)・[build](validation/item11-array-build.log)・[差分](validation/item11-array-diff.log)）。Goの3設定でfixtureの期待値を確認したが、これは公式PHPとの一致を意味しない。

配列要素も内部専用の`referenceCell`を保持する。`phpArray.get`は値を返し、参照を結ぶ操作だけが`getRaw`で束縛を取得する。通常代入・複合代入・増減は既存セルへ書き込み、`=&`は対象の要素だけを別セルへ付け替える。配列自体は従来の値コピーを保ち、**コピー後に**一方へ作った要素参照は他方へ現れない。以前の「**コピー前に**存在した要素参照は常に共有する」という記述は、上記の958件への修正で撤回した。[PHP公式マニュアルの参照と配列コピーの説明](https://www.php.net/manual/en/language.references.whatdo.php)と、修正節に示した複製処理を区別して参照する。要素や別名の`unset`は束縛だけを外し、残る別名を維持する。未定義キー・追記はnullの要素を作ってから参照し、入れ子のキーも評価する。通常読み取り、型判定、比較、foreach、引数展開では内部セルをPHP値として露出させない。

収録した境界は、変数↔要素・要素↔要素、付け替え、通常コピーの前後、外部別名をunsetしたコピー、要素のunsetと再作成、配列全体の置換、追記と入れ子、キー式の評価回数、キー変換時に警告ハンドラが投げた後の状態。単純な配列自己循環は、Go実装の制限として明示的な`gophp:`エラーで拒否し、PHP本来の診断とは扱わない。複雑な間接循環、プロパティ参照、参照引数・戻り値、参照捕捉、参照foreach、Traversableは未対応・未収録で、項目11全体は部分実装。今回の実装単位は**Go側の実装とローカル検証済み・公式比較待ち**であり、完了判定は3版の実比較後に行う。

[8.3対象](validation/item11-array-targeted-83-attempt.log)・[8.4対象](validation/item11-array-targeted-84-attempt.log)・[8.5対象](validation/item11-array-targeted-85-attempt.log)、[8.3全件](validation/item11-array-all-83-attempt.log)・[8.4全件](validation/item11-array-all-84-attempt.log)・[8.5全件](validation/item11-array-all-85-attempt.log)は、いずれもPodmanソケット接続拒否（公式PHP起動終了125）で比較器／`make compare`が終了2。各版0一致・0不一致・1環境エラー・955件未実施。fixtureの不一致は観測していない。次はMacでPodmanが利用できる状態になった後、以下を版ごとに順番に実行し、新規11件を対象比較して差異を直してから全955件を比較する。`make compare`は内部でローカルtest・buildも実行する。

```sh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  read -r suffix version target <<< "$spec"
  for fixture in tests/fixtures/success/reference_array_*.php; do
    log="docs/validation/item11-array-$(basename "$fixture" .php)-${suffix}-mac.log"
    printf 'command: compare_php.py --php php%s --php-version %s --fixtures %s --interpreter ./bin/gophp --php-version=%s\ntarget_php: %s\n' "$suffix" "$version" "$fixture" "$target" "$version" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php${suffix}" --php-version "$version" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=${target}" '{file}' 2>&1 | tee -a "$log"
    comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" | tee -a "$log"
  done
  log="docs/validation/item11-array-all-${suffix}-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\ntarget_php: %s\n' "$suffix" "$version" "$version" > "$log"
  make compare "PHP=./containers/php-version/bin/php${suffix}" "PHP_VERSION=${version}" 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" | tee -a "$log"
done
```

## 944件の公式比較（2026-10-18、今回の追加前の基準）

fixture・期待値は**944件**（正常524・実行時288・構文132）。[PHP 8.3.35](validation/reference-cells-944-83-mac.log)・[8.4.26](validation/reference-cells-944-84-mac.log)・[8.5.11](validation/reference-cells-944-85-mac.log)の版別生ログは、それぞれ944件の個別MATCH、**944一致・0不一致・0エラー**、`exit_code: 0`を記録する。新規の`reference_variables_*.php`7件も各版ですべてMATCH。収録fixtureの範囲で変数間の参照代入は公式PHP3版と比較済みで、追加7件の公式比較待ちは解消した。今回は保存済みログを確認し、比較は再実行していない。版別ログの結果を`make compare-all`の成功として扱わない。

実装範囲は変数間の`=&`に限る。配列要素・プロパティへの参照、参照引数・戻り値・捕捉・foreach等は未実装で、項目11全体は部分実装。下の接続拒否と暫定期待値の説明は、944件の公式比較**前**の履歴である。

## 実装直後の範囲（2026-10-18、944件比較前の履歴）

[移行判定](items6-10-integration-gate.md)の根拠となる版別生ログは[8.3](validation/warning-audit-final-83-mac.log)・[8.4](validation/warning-audit-final-84-mac.log)・[8.5](validation/warning-audit-final-85-mac.log)。各版937件の個別MATCH、937一致・0不一致・0エラー、`exit_code: 0`を確認した。`check.sh`のPASSはユーザー報告。項目6〜10の全仕様対応完了を意味しない。

今回、変数間の`$target =& $source`だけを追加した。fixture・期待値は**944件**（正常524・実行時288・構文132）で、既存937件は変更していない。新規7件の期待値は、[PHP公式マニュアルの参照の説明](https://www.php.net/manual/en/language.references.whatdo.php)と既存の値コピー・オブジェクト共有規則に基づく暫定値である。Goの3対象設定で新規7件は期待どおり動作するが、修正後の公式PHP比較は接続拒否で未実施。

## 格納と束縛

`variableScope`は変数名から格納値への表を維持する。参照される変数だけ、内部専用の`refKind`から`referenceCell`へのポインタを格納する。`readVariable`はセルの値を返し、`variableScope.slot`の通常書き込みはセルへ書く。`=&`は右辺のセルを作成・再利用し、左辺の**名前の束縛だけ**をそのセルへ付け替える。したがって通常代入は参照をコピーせず、配列は従来の値コピー、オブジェクトは従来の同一オブジェクト共有になる。`unset($name)`は名前だけを削除し、ほかの参照変数が使うセルを維持する。未定義の右辺変数はnullで作成する。自己参照は既存の値を保つ。

`global`／`static`は既存の`variableOwner`で読み書き先を決める。関数内で`global $x`とした名前を別の参照へ付け替えた場合は、ローカルの束縛を外して付け替え、グローバル変数自身の束縛を変えない。関数・メソッド呼び出しは値渡しの引数から新しいスコープを作る。`include`は呼び出し元スコープを使う。クロージャの値捕捉とトレースの引数表示では内部セルをPHP値として露出させない。

構文解析では`=`の直後の`&`だけを参照代入として扱い、ビット演算`&`、複合代入`&=`、型のintersectionと分離する。左辺・右辺が単純変数でない参照代入は、未選択分岐を含めて解析時に明示的に拒否する。これは**今回のGo実装の制限**であり、PHPが許す配列要素やプロパティへの参照の正規の診断を再現したものではない。

## 検証と残件

新規fixtureは参照代入式の値、双方向更新、通常代入の非共有、付け替え、複数参照と`unset`、未定義・自己参照、配列コピー、オブジェクト共有、`global`／`static`、読み込み元スコープを分けて観測する。Go単体テストは同一プロセスでの実行間分離と、未選択分岐に置いた未対応の参照元の早期拒否を確認する。[test](validation/item11-reference-test.log)・[vet](validation/item11-reference-vet.log)・[build](validation/item11-reference-build.log)は終了0。[対象試行8.3](validation/item11-reference-targeted-83.log)・[8.4](validation/item11-reference-targeted-84.log)・[8.5](validation/item11-reference-targeted-85.log)、[全件試行8.3](validation/item11-reference-all-83.log)・[8.4](validation/item11-reference-all-84.log)・[8.5](validation/item11-reference-all-85.log)は各版ともPodmanソケット接続拒否で終了2。fixtureの一致・不一致は各版0、環境エラー各1、修正後全944件と新規7件は公式比較未実施。ローカル期待値の成功を公式PHPとの一致と扱わない。

残件：参照引数・戻り値、クロージャ参照捕捉、参照foreach、配列要素・プロパティの参照、参照経由の型制約、Traversable展開、Iterator／Generator。次の単位は配列要素の参照と値コピーとの関係を、型付きプロパティや参照引数へ進む前に設計・検証する。項目11全体は部分実装。項目6〜10からのcallable、autoload、マジックメソッド、Reflection等の残件は[進捗文書](items6-10-progress.md)に維持する。

Macではリポジトリのルートで次を版ごとに実行する。`83 8.3.35 8.3`を`84 8.4.26 8.4`、`85 8.5.11 8.5`へ変え、ログも版ごとに分ける。新規7件はファイル名`reference_variables_*.php`を対象として個別に比較する。

```sh
set -o pipefail
for fixture in tests/fixtures/success/reference_variables_*.php; do
  log="docs/validation/item11-$(basename "$fixture" .php)-83-mac.log"
  printf 'command: compare_php.py --php php83 --php-version 8.3.35 --fixtures %s --interpreter ./bin/gophp --php-version=8.3\ntarget_php: 8.3.35\n' "$fixture" > "$log"
  python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php83 --php-version 8.3.35 --fixtures "$fixture" --interpreter ./bin/gophp --php-version=8.3 '{file}' 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
done
log=docs/validation/item11-reference-all-83-mac.log
printf 'command: make compare PHP=./containers/php-version/bin/php83 PHP_VERSION=8.3.35\ntarget_php: 8.3.35\n' > "$log"
make compare PHP=./containers/php-version/bin/php83 PHP_VERSION=8.3.35 2>&1 | tee -a "$log"
comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
```
