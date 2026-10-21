# 項目6〜10の統合検証と項目11への移行判定（2026-10-18更新）

## 項目11着手後の検証状況（2026-10-18、現在）

変数間の参照代入を加えた944件は、[項目11の版別生ログと対応範囲](item11-progress.md)で公式PHP3版の各版と**944一致・0不一致・0エラー、終了0**を確認した。項目11への移行判定は満たされたまま。項目6〜10の全仕様対応や項目11全体の完了を意味しない。以下の937件時点の移行判定と接続拒否は履歴である。

## 937件の公式比較と移行判定（2026-10-18、944件比較前の履歴）

[PHP 8.3.35](validation/warning-audit-final-83-mac.log)・[8.4.26](validation/warning-audit-final-84-mac.log)・[8.5.11](validation/warning-audit-final-85-mac.log)の版別生ログを確認した。各ログに937件の個別MATCH、**937一致・0不一致・0エラー**、`exit_code: 0`がある。`check.sh`の3版PASSはユーザー報告である。警告中断経路の今回定めた監査範囲と公式比較条件は満たしたため、**項目11へ移行可能**と判定する。項目6〜10の全仕様対応は完了していない。[項目11の進捗](item11-progress.md)と下記の過去の接続拒否記録を区別する。

## `define()`第3引数の追補（2026-10-18、937件比較前の履歴）

第3引数`null`は`bool`への引数変換時に非推奨診断が発生する。ハンドラが例外を投げた場合、Go側がそのerrorを捨てて定数を登録していたため、診断直後にerrorを返すよう修正した。第3引数`true`で出る警告は変換後の[PHP 8.4の`define`実装](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_builtin_functions.c)にある別経路で、登録処理に進む。この区別は公式ソースに基づく判断であり、今回のfixtureについて公式PHP3版の実比較結果ではない。新規`integration_warning_define_third.php`は、`null`でハンドラが投げた後の未登録と、`true`の警告をハンドラが処理した後の登録を対比する。fixture・期待値は**937件**（正常517・実行時288・構文132）。既存936件のソースと期待値は維持した。

[test](validation/warning-audit-third-null-test.log)・[vet](validation/warning-audit-third-null-vet.log)・[build](validation/warning-audit-third-null-build.log)は終了0。版別の対象比較[8.3](validation/warning-audit-third-null-targeted-83.log)・[8.4](validation/warning-audit-third-null-targeted-84.log)・[8.5](validation/warning-audit-third-null-targeted-85.log)、全件比較[8.3](validation/warning-audit-third-null-all-83.log)・[8.4](validation/warning-audit-third-null-all-84.log)・[8.5](validation/warning-audit-third-null-all-85.log)はすべてPodmanソケット接続拒否で終了2。各版0一致・0不一致・1環境エラーで、対象1件および修正後全937件は公式比較未実施。新規期待値は暫定で、935件の公式一致は変更前の履歴。

Macでの再比較はリポジトリのルートから版ごとに以下を実行する。`php83`／`8.3.35`／`8.3`をそれぞれ`php84`／`8.4.26`／`8.4`、`php85`／`8.5.11`／`8.5`に変えて残る2版も実行し、ログ名も版ごとに分ける。パイプライン直後の`$?`を記録する。

```sh
set -o pipefail
python3 skills/php-go-interpreter/scripts/compare_php.py --php ./containers/php-version/bin/php83 --php-version 8.3.35 --fixtures tests/fixtures/success/integration_warning_define_third.php --interpreter ./bin/gophp --php-version=8.3 '{file}' 2>&1 | tee docs/validation/warning-audit-third-null-targeted-83-mac.log
comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> docs/validation/warning-audit-third-null-targeted-83-mac.log
make compare PHP=./containers/php-version/bin/php83 PHP_VERSION=8.3.35 2>&1 | tee docs/validation/warning-audit-third-null-all-83-mac.log
comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> docs/validation/warning-audit-third-null-all-83-mac.log
```

修正後の3版一致が確認できるまで、項目11への移行判定は**条件付き**とする。

## 警告中断経路の追加監査（2026-10-18、936件時点の履歴）

比較済みの基準は935件（正常515・実行時288・構文132）で、[8.3](validation/items6-10-compound-final-compare-83-mac.log)・[8.4](validation/items6-10-compound-final-compare-84-mac.log)・[8.5](validation/items6-10-compound-final-compare-85-mac.log)の生ログが各935一致を記録する。今回、`internal/php`の`warn`／`diagnostic`呼び出しを棚卸しし、**定数登録や後続式の副作用に届く経路**に監査範囲を絞った。新規`integration_warning_define_state.php`を正常系へ1件追加し、fixture・期待値は**936件**（正常516・実行時288・構文132）。既存935件のソース・期待値は変更していない。

| 経路 | 確認したコードと既存の証拠 | 今回の判断 |
| --- | --- | --- |
| 読み込み・include_path、配列キー・動的プロパティ作成 | `include.go`／`include_path.go`、`aggregate.go`は診断から返るerrorを呼び出し元へ渡す。`error_handler_include_throw.php`、`include_path_handler_throw.php`、`integration_array_key_handler.php`、`integration_warning_transaction.php` | 既存935件の比較済み範囲を維持。動的プロパティが例外後も作成される条件は公式実測に従い、一律の巻き戻し規則を導入しない。 |
| 数値演算・ビット演算・型付きプロパティ | `arithmetic.go`／`bitwise.go`は整数変換のerrorを伝播する。`integration_modulo_handler.php`、`integration_compound_failure_state.php`、`integration_concat_handler_state.php` | 既存935件の比較済み範囲を維持。型付きプロパティの失敗後の格納値を演算・宣言型ごとに扱う。未収録のビット複合代入の型別境界へ一般化しない。 |
| `define(null, value)` | `constant_builtin.go`は第1引数の非推奨診断のerrorを捨て、ハンドラが投げても`defineConstant`へ進んでいた。既存fixtureにこの中断後の定数存在確認はなかった。 | errorを即時返すよう修正。新規fixtureでcatch後に空名の定数がないこと、別の正常な`define`が実行できることを確認。公式PHP3版比較待ち。 |
| `strlen(null)`・`count(..., null)`、オブジェクト数値cast、非配列`foreach` | `builtin.go`／`cast.go`／`foreach.go`には戻り値を直接使わない診断呼び出しがある。`eval`／`execute`の`pendingDiagnostic`による式・文境界の中断をコード上で確認した。 | 新規fixtureで、各ハンドラ例外後に右辺関数やループ本体のカウンタが増えないことをGo3設定で確認。公式PHP3版比較待ち。 |
| `construct.go`の静かな存在確認・unset、文字列添字、`aggregate_compare.go`、`closure.go`、`switch.go`、宣言時診断 | 一部は戻り値なしの診断と`pendingDiagnostic`に依存する。今回のコード監査では、診断直後に今回対象の定数登録と同じ外部状態の書き込みへ進む経路は確認しなかった。 | 構文・宣言の発生時点、参照・反復、未収録の添字変換は別の境界条件として残す。今回の公式一致や監査完了をそれらの全条件へ拡大しない。 |

当時の「第3引数の警告後も登録する」という説明は`true`指定時の警告に限られ、`null`の引数変換診断には当てはまらない。上の937件時点でこの判断を訂正した。

[ローカルtest](validation/warning-audit-test.log)・[vet](validation/warning-audit-vet.log)・[build](validation/warning-audit-build.log)は終了0。対象比較[8.3](validation/warning-audit-targeted-83.log)・[8.4](validation/warning-audit-targeted-84.log)・[8.5](validation/warning-audit-targeted-85.log)と版別全件[8.3](validation/warning-audit-all-83.log)・[8.4](validation/warning-audit-all-84.log)・[8.5](validation/warning-audit-all-85.log)は、Podmanソケット接続拒否で版確認時に停止し、比較コマンドは各終了2。fixtureの一致・不一致は各版0、環境エラー各1、新規1件と修正後全936件の公式比較は未実施。期待値は暫定。過去の935件一致を今回のコード変更後の成功として扱わない。

Macでの再比較は、ルートで次を版ごとに実行し、コマンド自身の終了コードを直後に記録する。`make compare`はMakefileの`PHP`と`PHP_VERSION`を指定し、Go側設定はそこから選ばれる。

```sh
cd /Users/kouji/dev/php
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  set -- $=spec
  log="docs/validation/warning-audit-targeted-$1-mac.log"
  printf 'command: compare_php.py --php php%s --php-version %s --fixtures tests/fixtures/success/integration_warning_define_state.php --interpreter ./bin/gophp --php-version=%s\ntarget_php: %s\n' "$1" "$2" "$3" "$2" > "$log"
  python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php$1" --php-version "$2" --fixtures tests/fixtures/success/integration_warning_define_state.php --interpreter ./bin/gophp "--php-version=$3" '{file}' 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
  log="docs/validation/warning-audit-all-$1-mac.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\ntarget_php: %s\n' "$1" "$2" "$2" > "$log"
  make compare "PHP=./containers/php-version/bin/php$1" "PHP_VERSION=$2" 2>&1 | tee -a "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
done
```

**移行判定：条件付き。** 上表の監査範囲で再現した定数登録の不具合はGo側で修正し、現時点で同じ範囲の未修正の既知不具合はない。ただし新規fixtureと全936件の公式PHP3版比較が未実施なので、項目11へ移行可能とはまだ判定しない。差異があれば修正して再比較し、3版一致を確認できれば移行する。項目11の最初の単位は参照セル基盤と参照代入（配列の値コピー・オブジェクト共有・参照束縛の区別）とする。参照foreach・Traversable等はその後に扱う。項目6〜10の全仕様対応は移行条件ではない。以下は935件までの比較と、さらに前の履歴である。

## 935件の公式比較と複合代入の訂正（2026-10-18、今回の変更前の履歴）

fixture・期待値は**935件**（正常515・実行時288・構文132）。版別の生ログは[PHP 8.3.35](validation/items6-10-compound-final-compare-83-mac.log)・[8.4.26](validation/items6-10-compound-final-compare-84-mac.log)・[8.5.11](validation/items6-10-compound-final-compare-85-mac.log)で、それぞれ935件のMATCH行と**935一致・0不一致・0エラー**を記録する。`integration_compound_failure_state.php`と追加の`integration_concat_handler_state.php`も3版でMATCH。版別`make compare`の終了コード0と`check.sh`のPASSはユーザー報告で、生ログには記載がない。`make compare-all`成功とは記録しない。

型付きプロパティの`/= 0`・`%= 0`・`+= []`・警告中断の`+= '2x'`は、収録例で例外後の格納値が0になる。通常変数・配列要素の警告中断では元の値を保つ。`.=`は左辺の値・宣言型によって異なる。追加fixtureでは、警告を投げるハンドラの下で文字列型プロパティは元の`a`を保ち、`int`・`float`・`?int`・`int|string`・`mixed`の代表例は0になる。したがって前回の「`.=`なら値は変わらない」という説明を撤回する。ビット複合代入には未検証の規則を広げない。

[PHP 8.4の演算子実装](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_operators.c)と[型付きプロパティへの複合代入](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_execute.c)を設計の参考とした。以前の「結果がないから格納しない」という推論は公式実比較より優先しない。[934件の修正前集計](validation/items6-10-compound-934-before-fix-summary-mac.log)は各版933一致・1不一致・0エラー、[933件一致集計](validation/items6-10-compound-933-after-fix-summary-mac.log)は各版933一致・0不一致・0エラー。両者は集計だけの過去記録である。修正直後の[対象比較](validation/items6-10-compound-recheck-targeted-84.log)・[全件比較](validation/items6-10-compound-recheck-all-84.log)の接続拒否も履歴として残す。935件の公式比較条件は満たしたが、下記の未収録警告中断経路の追加監査は未完了である。

以下は934件時点に用意したMac再比較コマンド（履歴）。新版でも再比較する場合は版ごとに実行し、個別fixture行・集計・終了コードを保存する：

```sh
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  set -- $=spec
  log="docs/validation/items6-10-compound-result-targeted-$1-mac.log"
  python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php$1" --php-version "$2" --fixtures tests/fixtures/success/integration_compound_failure_state.php --interpreter ./bin/gophp "--php-version=$3" '{file}' 2>&1 | tee "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
  log="docs/validation/items6-10-compound-result-compare-$1-mac.log"
  make compare "PHP=./containers/php-version/bin/php$1" "PHP_VERSION=$2" 2>&1 | tee "$log"
  comparison_rc=$?; printf 'exit_code: %s\n' "$comparison_rc" >> "$log"
done
```

## 剰余複合代入の追加修正と移行条件（2026-10-17、934件追加前の履歴）

警告ハンドラがfloat→int変換時に例外を投げる`%=`で、例外後に変数・配列要素・型付きプロパティが更新される既知不具合を修正した。`integerOperand`のエラーを`modulo`と共通の複合代入へ伝える。bitwiseも同じ整数変換のエラーを捨てない。`integration_modulo_handler.php`と`integration_modulo_normal.php`の2件を追加し、計**933件**（正常513・実行時288・構文132）。Goの3設定、`make test`、`make vet`、`make build`は通過した。

公式PHP3版の対象比較と版別全件比較は、いずれもPodman接続拒否（公式PHP版確認の終了125）でfixture実行前に停止した。各版**0一致・0不一致・1環境エラー・933件未実施**、比較コマンド終了2。新規期待値は暫定。[8.3対象](validation/items6-10-modulo-targeted-83.log)・[8.4対象](validation/items6-10-modulo-targeted-84.log)・[8.5対象](validation/items6-10-modulo-targeted-85.log)、[8.3全件](validation/items6-10-modulo-compare-83.log)・[8.4全件](validation/items6-10-modulo-compare-84.log)・[8.5全件](validation/items6-10-modulo-compare-85.log)。[PHPの剰余演算子](https://www.php.net/manual/en/language.operators.arithmetic.php)は整数変換を明記し、[throwの公式説明](https://www.php.net/manual/en/control-structures.throw.php)は例外後に通常の実行が中断することを明記する。ただし、今回のfixtureについて公式PHPの実出力は未確認である。

項目11への移行判定は引き続き**条件付き**。今回判明した状態破壊はGo側で修正したが、933件の公式PHP3版比較と残りの警告中断経路の監査が必要。以下の931件と新規7件の記述は今回の修正前の履歴である。

主基準はPHP 8.4.26、比較対象は8.3.35・8.5.11。追加前924件は[ユーザーの版別集計](validation/callable-924-after-fix-summary-mac.log)で各版924一致・差異0・エラー0。ただし集計は6行で、個別fixture行と終了コードは含まない。`check.sh`のPASSはユーザー報告であり、`make compare-all`の成功ではない。今回の実装・期待値は公式PHPとの再比較前は暫定とする。

## 統合検証の選定

| 範囲 | 既存の証拠 | 今回不足していた組み合わせと追加したfixture | 追加する理由 |
| --- | --- | --- | --- |
| A 読み込み×型・継承 | `include_interface_trait.php`、`special_include_variance.php`、`include_integration_parse_catch.php`、`include_integration_throw_once.php` | `integration_include_type_callable.php`、`integration_finally_callable_include.php` | 読み込み先のinterface・trait・クラス・関数が、呼び出し元のstrict_types、第一級callable、`static::class`、例外伝播と一緒に動くかを観測する。ParseErrorと_once履歴は既存例を再利用する。 |
| B callable×型・可視性 | `callable_scope.php`、`callable_strict_include.php`、`callable_late_binding.php`、`callable_builtin_named_order.php` | `integration_include_type_callable.php`、`integration_enum_callable.php` | privateメソッドから生成したClosureと読み込み先関数を、型検査・名前付き転送・enum型と組み合わせる。既存の単独可視性例は重複させない。 |
| C 型付き・readonly×例外・副作用 | `typed_property_caught_errors.php`、`readonly_type_failure_state.php`、`dynamic_property_receiver_change.php`、`readonly_promotion_composite.php` | `integration_warning_transaction.php`、`integration_promotion_failure_order.php` | 型変換警告やTypeErrorの後に元の格納値が残るか、複数昇格引数の式が順に評価され、型エラー時にコンストラクタ本体が走らないかを観測する。 |
| D 警告ハンドラ×操作 | `error_handler_deprecation_throw.php`、`error_handler_include_throw.php`、`error_handler_register_null_inside.php` | `integration_warning_transaction.php`、`integration_array_key_handler.php`、`integration_increment_handler.php` | `pendingDiagnostic`だけに頼ると、例外を返した後の格納が先に行われ得る。動的プロパティ、配列キーの型変換、増減で状態が変わらないことを確認する。 |
| E try／catch／finally×呼び出し・読み込み | `exception_finally_chain.php`、`include_integration_nested_finally.php`、`exception_catch_bindings.php` | `integration_finally_callable_include.php` | includeから返った後に作ったClosureが例外を投げても、呼び出し元のfinally、catch、後続処理の順序が保たれるかを確認する。 |

補助ファイルは`tests/fixtures/support/integration_*.inc`に置き、独立fixtureには数えない。追加前の924件は変更していない。当時追加した7件はすべて正常系fixtureで、捕捉した失敗と失敗後の状態を出力する。この段落の「追加」は931件時点を指す。現在は上記の935件生ログで収録範囲の比較を確認済み。

## 見つけた不具合と修正

`$box->extra = 7`の動的プロパティ作成時、非推奨診断のハンドラが例外を投げても、Go側はプロパティを作成していた。`locate`の格納前診断から返ったエラーを呼び出し元まで返し、格納を止めた。同じ問題を起こし得る配列キー変換、配列への自動変換、未定義の格納先読み取り、増減中の診断も、値を書き換える前にエラーを返す経路へ接続した。`integration_warning_transaction.php`と`integration_array_key_handler.php`はcatch後の値を観測する。`integration_increment_handler.php`はハンドラ例外後の文字列を観測する。

[公式マニュアルの動的プロパティ](https://www.php.net/manual/en/language.oop5.properties.php)、[例外とfinally](https://www.php.net/manual/en/language.exceptions.php)、[set_error_handler](https://www.php.net/manual/en/function.set-error-handler.php)を仕様根拠として参照した。資料確認は今回の公式PHP実比較の代わりではない。

## 残件の分類

| 分類 | 関連項目・残件 | 前提と移行への影響 | 次の検証方法 |
| --- | --- | --- | --- |
| 完了済み | 新規11件を含む全935件の公式PHP3版比較 | 版別生ログで各版935一致・0不一致・0エラーを確認。終了コード0はユーザー報告。公式比較を移行阻害条件から外す。 | 新しいfixtureを追加した場合のみ対象・全件を再比較する。 |
| A 移行前 | 項目10の未収録の警告中断経路 | 今回、格納を伴う代表経路は修正したが、`construct.go`などの診断呼び出しの戻り値と後続副作用は全経路を監査していない。状態破壊が再現したら既知不具合として移行前に修正する。 | 未収録の`warn`／`diagnostic`呼び出しを列挙し、ハンドラがthrowする最小例で格納値・設定を確認する。 |
| B 項目11 | 参照代入・引数・戻り値・Closure参照捕捉・型制約・参照foreach、Iterator／IteratorAggregate／Traversable／Generator、Traversable引数展開、`is_callable`出力参照引数 | 参照セルと反復プロトコルが前提。今回の統合fixtureはこれらを値渡しで代用しない。移行の妨げではなく項目11の主対象。 | まず参照セルと配列値コピー／オブジェクト共有の分離を設計し、小さい公式PHP比較例を作る。 |
| C 後続 | `Closure::bind`／`bindTo`／`call`、`__invoke`とマジックメソッド、clone、Reflection・属性、ストリームラッパー、Composer autoload、PHP 8.4プロパティフック・非対称可視性 | 項目6〜10の収録fixtureには不要。固定vendorには`Closure::bind`と`__invoke`の使用があるが、Goでのvendor実行到達は未確認。項目11の参照・反復基盤と切り離して範囲を決める。 | 固定vendorで実際に到達する経路を小さい再現へ切り出し、必要な組み込みAPIを確認する。 |

## 検証状態と移行判定

現在の総数は**935件**（正常515・実行時288・構文132）。版別生ログに各版935一致・0不一致・0エラーがある。ローカルの`make test`・`make vet`・`make build`・`git diff --check`の終了0、および版別比較の終了0と`check.sh`のPASSはユーザー報告であり、この3本の比較ログからは終了コードを確認できない。931件当時のPodman接続拒否、934件修正直後の接続拒否、暫定期待値の記述はそれぞれ当時の履歴である。

移行判定は引き続き**条件付き**。公式PHP3版での収録935件の比較という条件は満たした。一方、上表の項目10「未収録の警告中断経路」の追加監査は未完了である。`construct.go`などで`warn`／`diagnostic`の戻り値を受けた後に副作用が進む経路を調べ、見つかった不具合を修正・比較してから項目11へ移行する。項目6〜10のPHP全仕様対応を完了とする判定ではない。項目11の最初の単位は参照セル基盤と、値コピーする配列・共有するオブジェクト・参照束縛の切り分けとする。

Macでの再実行例（zsh、ログにコマンド・対象版・個別結果・集計・実際の終了コードを残す）：

```sh
cd /Users/kouji/dev/php
set -o pipefail
for spec in '83 8.3.35 8.3' '84 8.4.26 8.4' '85 8.5.11 8.5'; do
  set -- $=spec
  for fixture in tests/fixtures/success/integration_*.php; do
    log="docs/validation/items6-10-integration-targeted-$1-${fixture:t:r}.log"
    printf 'command: compare_php.py --php php%s --php-version %s --fixtures %s --interpreter gophp --php-version=%s\ntarget_php: %s\n' "$1" "$2" "$fixture" "$3" "$2" > "$log"
    python3 skills/php-go-interpreter/scripts/compare_php.py --php "./containers/php-version/bin/php$1" --php-version "$2" --fixtures "$fixture" --interpreter ./bin/gophp "--php-version=$3" '{file}' 2>&1 | tee -a "$log"
    result=$?; printf 'exit_code: %s\n' "$result" >> "$log"
  done
  log="docs/validation/items6-10-integration-compare-$1.log"
  printf 'command: make compare PHP=./containers/php-version/bin/php%s PHP_VERSION=%s\ntarget_php: %s\n' "$1" "$2" "$2" > "$log"
  make compare "PHP=./containers/php-version/bin/php$1" "PHP_VERSION=$2" 2>&1 | tee -a "$log"
  result=$?; printf 'exit_code: %s\n' "$result" >> "$log"
done
```
