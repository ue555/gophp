# 基本文法・名前空間の追加範囲

現在（2026-10-19）：型付き参照を含む997件は[公式PHP3版の生ログ](item11-progress.md)で各版一致・終了0。参照・反復と後続段階の入口を追加した現行**1035件**はGoローカル検証済み。追加38件と全件の公式PHP3版比較はPodman接続拒否で未実施、期待値は暫定。項目11と全6段階は部分実装。[全体進捗](full-stack-development.md)。

今回の追加前の比較済み基準：[PHP 8.3](validation/reference-property-974-83-mac.log)・[8.4](validation/reference-property-974-84-mac.log)・[8.5](validation/reference-property-974-85-mac.log)の版別生ログで各974件の個別MATCHと終了0を確認した。以下の974件「未比較」は当時の履歴。

追加前の比較済み基準：[版別生ログ](item11-progress.md)で公式PHP3版が各959件一致・終了0。以下の959件「未比較」は比較前の履歴。

959件比較前（履歴）：参照要素を含む配列の値コピーと、呼び出し元に残る参照の共有判定を加えた**959件**はGoローカル検証済み。公式PHP3版はPodman接続拒否で未比較、追加分の期待値は暫定。項目11全体は部分実装。[対応範囲と試行ログ](item11-progress.md)。

959件への追加前（履歴）：変数間参照を加えた**944件**は[版別生ログ](item11-progress.md)で公式PHP3版とも944一致・0不一致・0エラー、終了0。項目11全体は部分実装であり、収録外の参照形式は未対応。[対応範囲と残件](item11-progress.md)。以下は944件比較前の履歴。

944件比較前の状況（履歴、2026-10-18）：既存937件は[版別生ログ](items6-10-integration-gate.md)で公式PHP3版と一致、各比較終了0。項目11の変数間参照を加えた944件はGoローカル検証済みだが、公式比較は環境エラーで未実施。[項目11の範囲](item11-progress.md)。以下は937件以前の履歴。

937件比較前の検証状況（履歴、2026-10-18）：fixture・期待値は**937件**。`define()`第3引数`null`の診断中断をGo側で修正した。公式PHP3版の対象・全件比較はPodman接続拒否で未実施。[最新の監査と移行条件](items6-10-integration-gate.md)。以下は936件以前の履歴。

936件時点の検証状況（履歴、2026-10-18）：警告中断経路の追加監査でfixture・期待値は**936件**。`define(null, ...)`の診断後の定数登録をGo側で修正した。既存935件の[版別生ログ](items6-10-integration-gate.md)は修正前の証拠で、新規・全件の公式PHP3版比較はPodman接続拒否で未実施。[監査範囲と移行条件](items6-10-integration-gate.md)を参照。以下の935件以前の結果は各時点の履歴。

935件時点の検証状況（履歴、2026-10-18）：fixture・期待値は**935件**（正常515・実行時288・構文132）。公式PHP [8.3.35](validation/items6-10-compound-final-compare-83-mac.log)・[8.4.26](validation/items6-10-compound-final-compare-84-mac.log)・[8.5.11](validation/items6-10-compound-final-compare-85-mac.log)の版別生ログは各版935一致・0不一致・0エラーを記録する。ログには個別fixtureのMATCH行があり、版別比較の終了コード0と`check.sh`のPASSはユーザー報告である。公式比較という移行条件は満たしたが、未収録の警告中断経路の追加監査は残る。[統合検証一覧・移行判定](items6-10-integration-gate.md)を参照。項目6〜10全体は部分実装のままである。

931件時点の統合検証（履歴、2026-10-17）：新規7件を含めてfixtureは**931件**（正常511・実行時288・構文132）。Goローカル検証は通過し、警告ハンドラ例外後の格納を止める修正を加えた。当時はPodman接続拒否で公式PHP3版との比較ができず、新規期待値は暫定だった。直前の[924件版別一致](validation/callable-924-after-fix-summary-mac.log)は当時の変更前の証拠。[後続の比較結果と移行条件](items6-10-integration-gate.md)を参照。

統合検証前の追補（2026-10-17）：callable／Closure基本経路を含むfixtureは**924件**（正常504・実行時288・構文132）。ユーザーが版別に`make compare`を実行した[6行集計](validation/callable-924-after-fix-summary-mac.log)でPHP 8.3.35・8.4.26・8.5.11の各版は**924一致・0不一致・0エラー**。追加16件の公式比較待ちは収録範囲で解消した。`check.sh`のPASSはユーザー報告で、集計に個別fixture行・終了コードはない。今回は比較を再実行していない。[機能範囲と残件](items6-10-progress.md)。以下の908件と当時のPodman接続拒否は履歴。

今回の変更前の収録範囲は**908件**（正常490・実行時287・構文131）。PHP 8.3.35・8.4.26・8.5.11は[修正後の版別集計](validation/special-types-908-after-fix-summary-mac.log)で各908一致・0不一致・0エラー。特殊型の追加20件の比較待ちは収録範囲で解消した。この集計には個別fixture行とコマンド終了コードがなく、修正6件と`check.sh`のPASSはユーザー報告である。以下の「再比較待ち」などは各節の当時の履歴。[現在の残件](items6-10-progress.md)。

修正直後の記録（履歴）：分類は908件（正常490・実行時287・構文131）。特殊型の後方参照2件を失敗系へ移し、宣言診断を修正した。[修正前の908件集計](validation/special-types-908-before-fix-mac.log)は3版で不一致とテストエラーがあり、この時点では修正後の公式再比較待ちだった。[進捗](items6-10-progress.md)を参照。次段落の分類と「公式比較待ち」も、さらに前の追加時点を指す。

追加直後の記録（履歴）：収録範囲は908件（正常492・実行時285・構文131）。`self`／`parent`／`static`型とautoloadなしの前方型参照をGo側で部分実装した。従来の888件は[版別集計](validation/enum-888-after-fix-summary-mac.log)で3版それぞれ一致済み。新規20件は[比較試行](validation/special-types-compare-84-attempt.log)がPodman接続拒否で終了2となり、当時は公式PHP比較待ちだった。以下の878件などの記述も各時点の履歴として読む。[現状と残件](items6-10-progress.md)。

878件時点の履歴：収録範囲は**878件**（正常471・実行時278・構文129）。ユーザーが版別に`make compare`を実行した後の[集計ログ](validation/readonly-878-after-fix-summary-mac.log)で、PHP 8.3.35・8.4.26・8.5.11が各878一致・0不一致・0エラーと確認した。851件以降の追加27件は収録fixtureの範囲で公式比較済み。集計ログには個別fixture行とコマンド終了コードがない。詳細は[最新の進捗](items6-10-progress.md)と[比較記録](compatibility.md)を参照。以下の833件などの記述は各時点の履歴である。

833件時点の収録対象は833件（正常443・実行時273・構文117）。ユーザーの[修正後の公式PHP3版比較](validation/typed-members-833-retry-mac.log)は各版833一致・0不一致・0エラー、終了0。型付きプロパティ／定数の基本形と複合代入の格納値返却を含む新規16件の比較待ちは、収録fixtureの範囲で解消した。これは保存済みMacログの確認・保存であり、比較の再実行ではない。動的インスタンスプロパティアクセス、未登録型の解決等は残り、項目8と項目6〜10全体は未完了。

817件時点の記録（履歴）：ユーザーの[修正後の公式PHP3版比較](validation/void-never-817-after-fix-mac.log)は各版817一致・0不一致・0エラー、終了0。`void`・`never`の29件と`new \Exception`の1件を含む追加30件の比較待ちは、この時点で解消した。[修正前の差異](validation/void-never-817-before-fix-mac.log)は履歴として残す。

750件時点の記録（履歴）：[修正後の公式PHP3版比較](validation/union-diagnostic-750-after-fix-mac.log)は各版750一致・0不一致・0エラー、終了0。`?int|string`の構文エラー文言を含む[対象比較](validation/union-nullable-syntax-targeted-mac.log)も各版1件一致した。当時のunion診断の組み込み型順序は引数・戻り値・メソッド署名に共通適用済みだった。[修正前の不一致と残件](items6-10-progress.md)を参照。

## 718件比較と727件追加直後の記録（履歴）

修正前に比較済みの収録対象は718件（正常393・実行時250・構文75）。ユーザーのMacでの[修正後の公式PHP3版比較](validation/item7-dynamic-static-718-after-fix-mac.log)は各版718一致・0不一致・0エラー、終了0。追加11件の比較待ちは解消した。今回は保存済みログを確認・保存し、比較は再実行していない。[修正前の失敗](validation/item7-dynamic-static-718-before-fix-mac.log)は履歴として残す。直接`Class::{throw ...}()`をGoが受理する差は、この718件の比較時点では別の既知不具合だった。項目3〜5以降の部分実装と未対応範囲は[項目6〜10の進捗](items6-10-progress.md)と[現在の互換性](compatibility.md)を参照。項目7と項目6〜10全体は未完了で、Go製PHPでのLaravel起動も未対応。

2026-10-15追補（当時）：直接`Class::{throw ...}()`の事前拒否をGo側に追加し、fixtureは727件となった。この時点では追加9件の公式PHP3版比較がPodman接続拒否で未実施だった。後続の冒頭727件成功ログで解消した。現在の対応・残件は[進捗](items6-10-progress.md)を参照。

## 項目3〜5時点の記録

以下の件数・「未実施」・「未対応」は各段階の記録時点を指す。現在の状態は上記とリンク先を優先する。

2026-10-13更新: 項目3〜5の収録範囲は、提供された[527件の公式PHP比較ログ](validation/phase2-527-comparison-provided.log)で3版それぞれ527一致・差分0・エラー0と確認した。以下の古い「未実施」は、その時点の検証履歴として読む。依存解決と公式PHPでのLaravel/API Platform最小アプリ成功は[別記録](validation/laravel-api-evidence-2026-10-13.md)にあり、Go側のLaravel対応を意味しない。

主基準 PHP 8.4.26、追加比較 PHP 8.3.35 / 8.5.11。共通の構文解析・評価を使い、既存の `--php-version` で対象版を選ぶ。着手前466件を維持し、追加41件を含む507件を各版の公式PHPと比較した。ユーザー環境の比較ログでは3版とも507一致・差分0・エラー0。stdout・stderr・終了コードの厳密比較である。[検証記録](validation/phase2-local-validation.md)にログの確認結果を記載した。その後のnullsafeレビュー修正で11件を追加して総数518件とした。追加分を含む公式PHP比較は[レビュー修正の検証記録](validation/nullsafe-review.md)に示すとおり未実施である。

| 分類 | 今回の実装範囲 |
| --- | --- |
| 条件式 | `??`、`??=`、完全／省略三項。`??` と代入は右結合、連続する無括弧の完全三項は拒否。未選択式は評価しない |
| アクセス | `?->` のプロパティ・メソッド・添字チェーン。null なら残りのチェーン、添字、引数を評価しない。`isset`・`empty`・`??` は静かな読み取りを行う。代入先と`unset`対象としての使用は実行前に拒否 |
| 分岐 | 中括弧／`endswitch` 形式の `switch/case/default`、緩い比較、fall-through、`break/continue` のレベル処理。`match` は厳密比較、左から腕を評価し、未一致で `UnhandledMatchError` |
| ビット演算 | 整数変換を伴う `& \| ^ ~ << >>` と `&= \|= ^= <<= >>=`。両辺が文字列の `& \| ^` はバイト演算。負のシフト幅は `ArithmeticError` |
| 文字列 | 二重引用符内の `$name`、`$a[key]`、`$a[0]`、`$o->property`、`{$a[$key]}`、`{$o->property}`。変数参照をトークン化して通常の式評価へ渡す。heredoc は同じ展開処理、nowdoc は展開しない。終端ラベルの単語境界と共通インデントを処理する |
| 名前 | セミコロン式と括弧式の `namespace`、グローバル括弧、`use` / `use function` / `use const` とグループ形式、`\Foo` / `Foo\Bar` / `Foo` / `namespace\Foo`。クラスと関数は大文字小文字を区別せず、定数は区別する。未修飾関数・定数だけグローバルにフォールバック |
| 定数 | `const`、`define`、`defined`、`constant`、`__NAMESPACE__`、名前付きクラスの `Foo::class`。`Foo::class` はクラスの存在確認をしない。文字列の動的定数名には `use` を適用しない |

優先順位は [PHP の演算子表](https://www.php.net/manual/en/language.operators.precedence.php) の対象部分集合に従い、高い順に単項、乗除、加減、シフト、連結、大小、等価、`&`、`^`、`|`、`&&`、`||`、`??`、三項、代入。`??` と `??=` の quiet fetch は未定義変数・欠落した要素・プロパティで通常の未定義警告を出さない。添字式そのものの副作用は残り、`??=` は受け手・添字を読み書きで使い回す。今回の版差fixture `runtime/nullsafe_new_dereference.php` は、`new stdClass()?->missing` を8.3で構文エラー、8.4/8.5で実行可能として扱う。8.5の既存の特定警告・コンパイル時 traceも従来どおり切り替える。収録した診断は今回の公式PHP比較で一致した。未収録の版差は未確認である。

今回の文字列展開は上表の形式に限定する。`${...}` は PHP で有効な旧形式であっても、今回の Go 実装では明示的に拒否する。メソッド呼び出しを伴う複雑な中括弧展開、可変変数、入れ子展開、Unicodeコードポイントや255超の8進エスケープも明示的に拒否する。単純展開に後続する複雑な連鎖は文字列として解釈されるPHP規則があるため、これらは現時点で個別検証が必要である。heredoc / nowdoc の CR 単独改行やすべてのインデント例は未確認。

項目3〜5時点の`const` 初期化で Go が受け付けるのは対応済みのリテラル・定数名・配列、およびそれらの単項・二項・三項演算からなる式である。関数呼び出し、`new`、文字列展開を含む式は拒否する。PHP の全定数式を実装した意味ではない。当時はクラス定数、`self::class` / `parent::class` / `static::class`、動的な関数・メソッド・クラス名、静的メンバー呼び出しは未対応だった。現在はクラス定数と相対クラス名・staticメンバーの部分集合を追加したが、動的な名前全般と型付き定数は残る。未収録の `match` のエッジケース、型変換や警告の版差にも追加の実比較が必要である。

[実行時の設計](phase2-design.md)、[依存コード調査](dependency-audit.md)、[調査アプリ](../research/laravel-api/README.md) を参照。Laravel と API Platform の起動、JSON HTTP、入力検証、SQLite CRUD は本段階の実装範囲ではない。

## nullsafeメソッドの受け手（2026-10-12追記）

`isset($o?->secret)`、`empty($o?->secret)`、`$o?->secret ?? "fallback"` の末端プロパティは静かに読み取る。一方、`$o?->secret?->take() ?? "fallback"` の `secret` はメソッド呼び出しの受け手なので、通常のプロパティ読み取りと同じアクセス権検査・警告を行う。先行する `?->` がnullで短絡した場合は、後続のプロパティ、添字、引数を評価しない。ルートの未定義変数を使った `$undefined?->take() ?? "fallback"` は警告なしで `fallback` となる。以前の「警告が必要」という指摘は、3版の実比較により訂正された。

この修正で9件を追加し、全527件となった。修正前に提供された[公式PHP実比較ログ](validation/nullsafe-receiver-before.log)では、privateプロパティのケースが3版とも不一致、未定義変数のケースは3版とも一致している。新しい関連ケースの期待値には、公式PHPで未比較の暫定値が含まれる。今回の環境ではPodmanソケットへの接続が拒否され、修正後の公式PHP比較は3版とも版確認で停止した。件数と再実行手順は[検証記録](validation/nullsafe-receiver-review.md)に示す。

その後、[527件の実比較ログ](validation/phase2-527-comparison-provided.log)の実体を受領し、3版すべて527一致・差分0・エラー0を確認した。以前のエージェント環境での未実施記録は当時の履歴として残す。Laravel/API Platformの依存解決と公式PHP HTTP実行は[別記録](validation/laravel-api-evidence-2026-10-13.md)で確認したが、Goでの実行は未対応である。

2026-10-13に着手した項目6〜10では、さらに11 fixtureを加えて538件とした。追加分はGoの期待値テストのみ通過し、公式PHP3版での比較はPodman接続拒否により未実施である。`include`、単純な `extends` / `instanceof`、`global`・`static` ローカル、変数関数の部分実装と残件は[進捗文書](items6-10-progress.md)に記録した。この段落以降も上記527件の実比較と新規11件の未比較を混同しない。
