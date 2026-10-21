# 実行時構造と今後の読み込み設計

現在（2026-10-19）：参照セルの型付きプロパティ制約に加え、クロージャの参照捕捉・既存変数の参照引数・配列変数の参照foreachを同じセルへ接続した。Iterator／IteratorAggregate／Countableは既存のクラス・メソッド呼び出しを利用する。後続段階の匿名クラスは同じクラス宣言・生成経路へ、分解代入は既存格納先へ接続した。**997件**は公式PHP3版と各版一致・終了0、現行**1035件**の追加38件はGoローカル検証済みで公式比較待ち。[設計・未対応境界](full-stack-development.md)。

今回の追加前の比較済み基準：[PHP 8.3](validation/reference-property-974-83-mac.log)・[8.4](validation/reference-property-974-84-mac.log)・[8.5](validation/reference-property-974-85-mac.log)の版別生ログで各974件の個別MATCHと終了0を確認した。以下の974件「未比較」は当時の履歴。

今回の追加前の比較済み基準：959件は[版別生ログ](item11-progress.md)で公式PHP3版が各959一致・終了0。以下の959件「未比較」は比較前の履歴。

959件比較前（履歴）：関数呼び出し時に退避したローカルスコープも、配列要素参照の共有判定へ含める。別の格納先が保持しない要素参照はコピー先で通常値へ分離する。**959件**はGoローカル検証済み、公式PHP3版は環境エラーで未比較。[項目11の設計と境界](item11-progress.md)。

今回の修正前（履歴）：配列書き込み時、別の格納先が保持しない要素参照はコピー先で通常値へ分離する。**958件**はGoローカル検証済み、公式PHP3版は環境エラーで未比較。[項目11の設計と境界](item11-progress.md)。

958件への修正前（履歴、2026-10-18）：配列要素も内部参照セルへ接続し、通常の値コピーは配列パスを分離、コピー前に存在した要素セルだけを共有する。**955件**はGoローカル検証済み、公式PHP3版はPodman接続拒否で未比較。設計と制限は[項目11の記録](item11-progress.md)。

955件追加前の基準（履歴）：変数間参照のセルと束縛を加えた**944件**は[版別生ログ](item11-progress.md)で公式PHP3版とも944一致・0不一致・0エラー、終了0。参照引数・配列要素参照などへ拡張する設計と残件は[項目11の記録](item11-progress.md)を参照。以下は944件比較前の履歴。

944件比較前の状況（履歴、2026-10-18）：項目11では既存`variableScope`の変数名束縛に共有参照セルを加え、通常の読み書きはセルを透過的に参照する。詳細は[項目11の設計](item11-progress.md)。既存937件の公式PHP3版比較は[版別生ログ](items6-10-integration-gate.md)で全件一致・終了0。参照追加後の944件は公式比較待ち。以下は各時点の履歴。

937件比較前の警告中断監査（履歴、2026-10-18）：`define()`第3引数`null`の`bool`変換中に診断ハンドラが投げたら、`constant_builtin.go`でerrorを返して登録前に中断する。第3引数`true`の明示的な警告と登録は別経路。fixture・期待値は937件でGoローカル検証済み、公式PHP3版比較は環境エラーで未実施。[棚卸しと移行条件](items6-10-integration-gate.md)。以下は936件以前の履歴。

936件時点の警告中断監査（履歴、2026-10-18）：`define(null, ...)`は第1引数の非推奨診断が例外になった時点で定数登録へ進まないよう、`constant_builtin.go`でerrorを返す。`pendingDiagnostic`は他の値返却経路の式・文境界で使われるが、外部状態を変える`define`では診断直後に中断する必要があった。新規fixtureを含む936件はGoローカル検証済み、公式PHP3版比較待ち。[棚卸しと移行判定](items6-10-integration-gate.md)。以下の935件の比較結果は修正前の履歴。

複合代入失敗結果の追補（2026-10-18、936件追加前の履歴）：[935件の版別生ログ](items6-10-integration-gate.md)で、型付きプロパティの`/= 0`・`%= 0`・`+= []`・警告中断の`+= '2x'`による値0を含む全fixtureが公式PHP3版と一致した。`evalCompound`では確認済みの演算・被演算子・エラーを限定して0を格納する。通常変数・配列要素には適用せず、ビット演算にも拡大しない。`.=`も左辺の宣言型で異なり、文字列型の値保持を全型へ一般化しない。以前の「結果のない演算は格納しない」という判断、および`.=`の一律な値保持という説明は訂正した。比較コマンド終了0はユーザー報告で生ログにはない。次段落の933件は過去の記録。

剰余複合代入の追補（履歴、2026-10-17）：`integerOperand`は整数・適合フラグ・errorを返し、floatからintへの変換警告でハンドラが投げた例外を`modulo`・`evalCompound`経由で格納前に伝える。bitwiseも共通変換のエラーを伝播する。現行933件はGoローカル検証済み、公式PHP3版比較はPodman接続拒否で未実施。[移行条件](items6-10-integration-gate.md)。次段落の931件は修正前の履歴。

統合検証の前段階（履歴、2026-10-17）：警告ハンドラから返った例外を`locate`、`arrayKey`、`incrementValue`、数値演算の格納前まで伝える。`pendingDiagnostic`による式末尾の中断だけでは、例外後に書き込みが済む経路を防げなかった。新規7件を含む現在の931件はGo側で検証済みで、公式PHP3版との比較はPodman接続拒否により未実施。[統合検証一覧と移行条件](items6-10-integration-gate.md)。下の924件一致は変更前の履歴。

統合検証前の検証履歴（2026-10-17）：callable／Closureの基本経路を含む収録924件は、ユーザーの[版別6行集計](validation/callable-924-after-fix-summary-mac.log)で公式PHP3版それぞれ924一致・0不一致・0エラー。追加16件の比較待ちは収録範囲で解消した。集計に個別fixture行と終了コードはなく、今回は保存済みログの確認のみで比較を再実行していない。下記の接続拒否・比較待ちは当時の履歴。[進捗](items6-10-progress.md)。

今回の変更前の検証（2026-10-16）：相対型と後方型宣言の修正を含む収録908件は、[ユーザーの版別比較集計](validation/special-types-908-after-fix-summary-mac.log)でPHP 8.3.35・8.4.26・8.5.11の各版と一致した。追加20件の比較待ちは収録範囲で解消した。集計に個別fixture行・終了コードはなく、修正6件と`check.sh`のPASSはユーザー報告である。以下の「比較待ち」は修正直後の履歴。[進捗と残件](items6-10-progress.md)。

## 後方型宣言の継承検査と診断修正（2026-10-16）

[修正前の908件比較](validation/special-types-908-before-fix-mac.log)で、継承時の戻り値型を比較するために後方のクラスを先行登録する実装が公式PHPと異なることが分かった。`extends`の依存だけは既存の前方登録で解決し、署名内の型名は比較のためだけに登録しない。先行する実行文がある場合の宣言到達時までの保留は維持する。署名診断では`self`／`parent`を宣言元へ展開し、未登録型が両側にある戻り値比較では子側から利用不能な型名を報告する。親なし`parent`には引用符を付ける。修正後の公式比較は待機中で、[進捗](items6-10-progress.md)の期待診断は暫定である。

## 相対型と型情報不足の宣言検査（2026-10-16、追加時点の履歴）

型文字列は既存のunion選択肢・intersection構成型へ解析し、`relativeTypeName()`が各構成型の`self`／`parent`を宣言元へ、実行時の`static`をcalled scopeへ解決する。メソッドの継承比較は双方の`methodDecl.owner`を渡し、弱い型変換から独立した`typeAcceptsScoped()`を使う。`static`は戻り値でlate-boundのまま比較する。型付きプロパティと定数の継承比較も宣言元を渡す。未知のクラス名は同一名なら比較できるが、階層関係を要する場合は関係を推測しない。診断では、実際に失敗した型関係にある未登録型だけを調べる。

`executeProgram()`の事前登録は同じファイルの無条件クラス定義を利用する。先行する実行文が未登録の型を定義し得る場合は、子宣言の登録を到達時まで保留する。includeを先読みしたり、未選択分岐のクラスを登録したりしない。これはautoload一般対応ではなく、循環・未登録型の宣言時点、特殊型の全構文は[進捗](items6-10-progress.md)に残す。新規20件の期待値は、この時点ではPodman接続拒否により暫定だった。[版別比較試行](validation/special-types-compare-84-attempt.log)。

## enum予約メソッドの検査時点とtrait由来の識別（2026-10-16）

直接宣言された`cases()`およびbacked enumの`from()`／`tryFrom()`は、構文解析中に拒否せず、`classDecl`のメソッド表へ保持して宣言登録時に検査する。先行する実行文がある場合は該当enumの事前登録を保留し、宣言時点の出力順とenum宣言行の診断を維持する。trait合成後の`methodDecl.traitSource`を使い、trait由来の`from()`／`tryFrom()`を直接宣言と区別する。プロパティ宣言禁止は解析時の検査を維持する。現在の[888件の版別比較集計](validation/enum-888-after-fix-summary-mac.log)は3版で全件一致を記録する。ただし個別fixture行と終了コードは含まれない。[進捗](items6-10-progress.md)を参照。

## 通常enumのユーザー定義from／tryFrom（2026-10-16、修正前の履歴）

`classDeclaration`の直接メソッド検査と`composeTraits`後の検査の両方で、`from`／`tryFrom`を禁止する条件にbacking型の有無を加えた。通常enumでは`callEnumStatic`が未処理として既存のユーザー定義staticメソッドへ渡す。`cases`は両enumで引き続き禁止する。例示コードで使う戻り値`self`／`?self`はメソッド呼び出し中の宣言元スコープで検査する。この最小対応を特殊型全般の完成とは扱わない。[進捗](items6-10-progress.md)を参照。

## enumの実行構造（2026-10-16、実装時点の記録）

`classDecl`へbacking型と宣言順のcase一覧を追加し、登録時にcaseごとの単一オブジェクトを作ってclass定数表に保持する。`phpObject.enumCase`で通常オブジェクトと区別し、書き込み・unsetを拒否する。case同士の緩い比較もプロパティ値ではなくオブジェクト同一性で判定する。`UnitEnum`と`BackedEnum`は組み込みinterfaceとして実行単位のクラス表へ登録する。caseオブジェクトの複製はしない。`cases()`と`from()`／`tryFrom()`は専用経路のため、既存の引数束縛・トレースとの完全統合と厳密な診断が残る。実装直後の[接続拒否ログ](validation/enum-compare-all-attempt.log)は当時の履歴であり、後続の[888件の版別比較集計](validation/enum-888-after-fix-summary-mac.log)で収録fixtureの一致を確認した。[進捗](items6-10-progress.md)を参照。

仕様調査には[PHP公式enum RFC（Implemented）](https://wiki.php.net/rfc/enumerations)、[公式マニュアルのbacked enum](https://www.php.net/manual/en/language.enumerations.backed.php)、[通常オブジェクトとの差](https://www.php.net/manual/en/language.enumerations.object-differences.php)、[PHP 8.4のenum実装](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_enum.c)を用いた。PHP実装ではbacked値の重複をbacked table構築時に検出するため、Goでも宣言時に重複候補を記録し、case取得または組み込みAPI利用時に報告する。ただしこの遅延解決の全境界と診断は未比較である。これらは仕様根拠であり、新規fixtureを公式PHP3版で実行した証拠ではない。

## readonlyとプロパティ昇格（2026-10-16、収録878件は版別の集計で一致）

`parameter`は昇格時の可視性・readonlyを保持し、コンストラクタの解析時に対応する`propertyDecl`を追加する。昇格引数の既定値は引数にのみ残し、プロパティの宣言時初期値にはしない。呼び出し時は既存の名前付き／位置引数束縛と型検査を終えてから、コンストラクタ本体の前にプロパティへ格納する。親から継承したコンストラクタとtraitから取り込んだコンストラクタも宣言元のスロットを使う。

readonlyの初期化済み判定は、既存のオブジェクト別プロパティ格納にキーがあるかどうかで行い、`null`と未初期化を混同しない。書き込みは`storeObjectProperty`へ集約し、初期化可能スコープ、再書き込み、型変換を順に判定する。readonly classは宣言のフラグで、全プロパティの実効readonlyと動的プロパティ作成禁止に使う。PHP 8.3の宣言元限定、8.4以降の子クラスへの初期化権限拡大は版設定で分ける。配列は既存の値コピー経路で更新を拒否し、格納したオブジェクトの内部変更は許す。clone・参照・フック・明示的なset可視性は後続課題。[公式プロパティ説明](https://www.php.net/manual/en/language.oop5.properties.php)、[昇格RFC](https://wiki.php.net/rfc/constructor_promotion)、[検証状況](items6-10-progress.md)を参照。

継承時は親子のreadonlyフラグが異なれば宣言エラーにする。通常の子がreadonly親を継承する方向を許す案は[readonly amendments RFCの提案1](https://wiki.php.net/rfc/readonly_amendments)で否決された。既存の逆方向の検査と併せて`registerClass`で親の解決後に判定し、readonly親子の合法な継承は維持する。収録fixtureの[版別比較後の集計](validation/readonly-878-after-fix-summary-mac.log)は3版すべて878件一致した。

trait合成後の宣言検査では、readonly classへ取り込む非readonlyプロパティを元traitの名前を使って拒否する。子メソッドから親のprivateプロパティ名へ書き込む場合は親privateスロットへアクセスせず、別の公開された動的プロパティを作る。作成時の非推奨診断は宣言された親privateプロパティの存在で抑止しない。両ケースの修正前[公式PHP比較](validation/readonly-878-before-fix-mac.log)と、修正後の[集計結果](validation/readonly-878-after-fix-summary-mac.log)を区別する。後者は個別fixture行・終了コードを含まず、両ケースの解消はユーザー報告である。


動的インスタンスプロパティ名の追加設計（2026-10-16）：`propertyExpr`と`nullsafeExpr`に名前式を保持する。読み取り・静かな確認・書き込み・unsetでは名前式を一度評価してから、既存の`propertyKey`、可視性検査、宣言元別privateスロット、型付きプロパティの格納処理へ接続する。通常代入では名前式を準備時に一度評価し、変数で指定した受け手は右辺評価後の書き込み時に解決する。これは[公式PHPの比較ログ](validation/dynamic-property-851-before-fix-mac.log)で、名前式または右辺が受け手変数を書き換えた場合の格納先から確認した。`??=`は静かな読み取り後に同じ受け手と名前を使う。nullsafeで受け手がnullなら名前式と後続添字を評価しない。動的名からのメソッド呼び出しは今回の対象外として明示的に拒否する。当初はこちらのPodman接続拒否で比較できなかったが、後続のMac実行で8.3・8.4の全851件と、別実行による8.5の全851件の一致を確認した。[進捗と比較手順](items6-10-progress.md)を参照。

型付きメンバーの追加設計（2026-10-16）：`propertyDecl`と`classConstantDecl`が型情報を保持する。インスタンスの型付きプロパティは未初期化時にスロットを持たず、同じ名前の親子privateプロパティは宣言元を含む既存キーで分離する。staticは共有スロットに`initialized`を持ち、親子の再宣言は新スロットを作る。読み取りと書き込みは異なる経路とし、書き込みは共通の型変換を経てから格納する。型付き定数は値を宣言時に検査し、代入用の弱い型変換を適用しない。収録833件の[公式PHP3版比較](validation/typed-members-833-retry-mac.log)は一致した。動的インスタンスプロパティアクセスなどの残件は[進捗](items6-10-progress.md)を参照。

2026-10-15の診断修正：nullableな`void`／`never`は構文解析の失敗としてではなく、宣言検査の`fatalDiagnostic`として扱う。`void`の文言は先頭を`Void`とし、`never`は小文字のままにする。コンストラクタの戻り値型は`Method Class::__construct() cannot declare a return type`とする。これらは[公式PHP3版の不一致ログ](validation/void-never-817-before-fix-mac.log)に基づき修正し、[修正後の全件比較](validation/void-never-817-after-fix-mac.log)で収録817件が3版すべて一致した。

## 名前空間内の完全修飾クラス生成（2026-10-15、収録fixtureは公式比較済み）

字句トークンに開始オフセットを保持し、式の先頭の`new`と後続の`\ClassName`が空白・改行・コメントで離れている場合は、一つの修飾名へ連結しない。`new`を構文キーワードとして扱い、クラス名だけを既存の`readName()`／`resolveName()`で解決する。これにより名前空間内の`new \Exception`はグローバルの`Exception`へ解決される。連続した名前トークンの扱いは維持する。[回帰fixture](../tests/fixtures/success/namespace_new_fqcn.php)は[公式PHP3版の全件比較](validation/void-never-817-after-fix-mac.log)で一致した。

## void／never戻り値制御（2026-10-15、収録29件は公式比較済み）

`methodDecl.returnType`の`void`と`never`は値型への変換に渡さない。`validateCallable()`が各callableの本文を入れ子宣言の境界で区切って走査し、`void`の値付きreturnと`never`の明示的returnを実行前に拒否する。`invokeWithCapture()`は`void`の裸のreturnと末尾到達をnullの呼び出し結果へ変換し、`never`の末尾到達を`TypeError`として呼び出し元へ伝播させる。既存の`try`／`catch`／`finally`の中断経路を使い、throwを値として処理しない。継承署名の`typeAccepts()`では`never`を最小戻り値型、`void`を通常の値型と異なる型として扱う。型なしの親戻り値は子の`void`を受け入れる。メソッド診断の8.3と8.4以降の文言差は宣言に記録した対象版で判断する。[公式PHP3版の比較結果](validation/void-never-817-after-fix-mac.log)。

## 宣言時の数値型と実行時変換の分離（2026-10-15）

継承署名の`typeTermAccepts()`では`int`を`float`の下位型として扱わない。戻り値の共変・引数の反変・union内の包含判定に共通の規則を適用する。実行時の`coerceType()`では整数をfloatに変換でき、`strict_types=1`でもその動作を維持する。[PHP 8.4の継承検査](https://raw.githubusercontent.com/php/php-src/PHP-8.4/Zend/zend_inheritance.c)を根拠とした。追加直後は接続拒否で公式比較できず期待値を暫定としたが、後続の[787件成功ログ](validation/type-variance-787-compare-all-mac.log)で収録fixtureの3版一致を確認した。

## 複合型の継承署名包含（2026-10-15、収録fixtureは公式比較済み）

`types.go`の`typeAccepts(receiver, source)`は、宣言型をunionの選択肢、各選択肢をintersectionの構成型として比較する。`source`の各選択肢が`receiver`の少なくとも一つに含まれる必要があり、intersectionの各要求型は、元のintersectionのいずれかの構成型から導ける必要がある。名前の同一性に加え、登録済みクラスの親連鎖とinterface継承・実装を使う。引数では子の型をreceiver、親の型をsourceとし、戻り値では逆向きにする。型なしと`mixed`は許容値集合の最上位だが、署名表示は区別する。`?T`は`null`を含む集合として扱う。弱い型変換は実行時の`coerceType()`だけに残す。

この判定はPHPの宣言時DNF冗長性検査と別であり、クラスを読み込んでまで包含を証明しない。未登録の名前付き型の階層関係、`self`／`parent`／`static`、`void`／`never`、複合型の未収録境界は残る。新規fixtureは[進捗](items6-10-progress.md)に記録し、[787件成功ログ](validation/type-variance-787-compare-all-mac.log)で収録範囲の公式PHP3版一致を確認した。[PHP公式の変性規則](https://www.php.net/manual/en/language.oop5.variance.php)と[DNF RFC](https://wiki.php.net/rfc/dnf_types)を設計の参照とした。

## intersection／DNF型の表現（2026-10-15追加直後の履歴）

`declaredType`はunionの各選択肢を`typeAlternative`として保持し、各選択肢内の`types`をintersectionの構成型とする。単純型・nullable・unionも同じ形で判定し、型表示は`declaredTypeText()`へ集約する。名前解決は宣言の解析時に行い、実行時の照合はクラス／interfaceの関係を各構成型に適用する。DNFのスカラー選択肢では既存unionの弱い変換順序を使う。構成型の大小文字と順序を正規化したキーで、同じ型集合の継承署名を比較する。[PHP型宣言の公式説明](https://www.php.net/manual/en/language.types.declarations.php)と[採用済みDNF RFC](https://wiki.php.net/rfc/dnf_types)を構文・重複判定の参照にした。クラス階層による複合型の完全な分散検査は次段階に残す。[今回の暫定検証と再比較手順](items6-10-progress.md)。

2026-10-15追補：DNFの2つのintersection選択肢は、片方の構成型集合が他方の真部分集合であれば、構成型の多い方をコンパイル時に冗長として拒否する。名前解決後の大小文字を無視した集合比較のみを使い、クラス階層・autoloadは参照しない。選択肢と構成型の順序を入れ替えても同じ判定になる。新規4件は追加直後には公式比較待ちだったが、後続の[769件成功ログ](validation/dnf-769-after-fix-mac.log)で収録範囲の一致を確認した。

## union型の診断表示（2026-10-15、収録fixtureは公式比較済み）

型の適合判定では宣言された構成型を保持する一方、診断には`declaredTypeText()`を共通利用する。クラス名は宣言順を保ち、組み込み型はPHPの`zend_type_to_string_resolved()`の固定順で表示する。これにより引数・戻り値・継承メソッド署名・既定値検査で表示順が揃う。[PHP 8.4のソース](https://sources.debian.org/src/php8.4/8.4.26-1/Zend/zend_compile.c/)を参照。追加直後はGoローカルの750件だけが通過し、Podman接続拒否で公式PHP3版の比較は未実施、期待診断は暫定だった。[検証記録](items6-10-progress.md)を参照。

後続の[修正前実比較](validation/union-diagnostic-750-before-fix-mac.log)では、3版とも749件一致し、`?int|string`の構文エラー文言だけが不一致だった。`typeDeclaration()`は`?int`の直後の`|`を、次のunion構成型として読まずに構文エラーにする。この位置の期待トークンは引数宣言では変数であり、実行前に拒否する。その後の[Macでの修正後比較](validation/union-diagnostic-750-after-fix-mac.log)は3版とも750件一致し、[対象比較](validation/union-nullable-syntax-targeted-mac.log)でも各版1件一致した。これは保存済みログの確認であり、比較の再実行ではない。[現状](items6-10-progress.md)を参照。

## 動的staticメソッド名の直接`throw`（2026-10-15、727件公式比較済み）

[修正後の成功ログ](validation/static-throw-727-compare-all-mac.log)でPHP 8.3.35・8.4.26・8.5.11の各727件が一致し、追加9件の比較待ちは解消した。今回は保存済みログを確認・保存し、比較は再実行していない。

`staticCallExpr`の`nameExpr`が直接`throwExpr`となる場合だけ、`validateStatements`の再帰走査で実行前にコンパイル致命エラーを返す。括弧付きの直接形も同じASTとなる。三項式内の`throw`や関数呼び出しは別の評価結果を持つため一律に禁止しない。PHPの[`zend_compile_static_call`](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_compile.c)はコンパイル済みメソッド名が定数の非文字列なら`E_COMPILE_ERROR`とし、`zend_compile_throw`は式の結果として真の定数を返す。主ファイルは実行前の検査、読み込み先は読み込み中の検査に乗せ、既存の致命的な宣言エラー経路を維持する。追加9件のGoローカル検証は成功した。この時点ではPodman接続拒否により公式PHP比較が未実施だったが、後続の上記成功ログで解消した。[進捗と再比較手順](items6-10-progress.md)を参照。

## 動的staticメンバー（2026-10-15、718件公式比較済み）

ユーザーのMacでの[修正後全件比較](validation/item7-dynamic-static-718-after-fix-mac.log)はPHP 8.3.35・8.4.26・8.5.11で各718一致・0不一致・0エラー、終了0。`class_dynamic_static_order.php`も3版で一致した。今回は保存済みログの確認・保存であり、比較は再実行していない。以下の修正前の差異は履歴である。

static呼び出し・プロパティASTに、ソース上で解決済みの固定名とは別にクラス式とメンバー名式を保持する。`Class::$name`の裸の読み書きは固定プロパティ`name`を指し、`Class::$name()`では変数`$name`の値をメソッド名として使う。`Class::{expr}()`と`Class::${expr}`は式を評価し、動的クラス式も評価した後、既存のクラス表・可視性・遅延静的束縛へ渡す。staticプロパティの書込み準備時にはクラス式と名前式を一度評価して固定し、右辺の評価と書込みで再実行しない。`isset`・`empty`・`??`の静かな取得も同じ解決を使う。`self::`・`parent::`・`static::`の転送呼び出しは従来のcalled scopeを保持する。文字列による動的クラス指定は名前空間インポートや`self`等のスコープキーワードとして再解釈せず、クラス表を直接引く。[PHP公式のstaticメンバー規則](https://www.php.net/manual/en/language.oop5.static.php)を参照。新規11件の修正前比較では10件が公式PHP3版で一致し、残る1件は成功fixtureの前提違反による検証エラーだった。式形式・診断の未収録境界と動的クラス定数は残件。

修正前の[Mac全件比較](validation/item7-dynamic-static-718-before-fix-mac.log)では、メソッド名式に直接`throw`を書いた成功fixtureが公式PHP3版とも実行前Fatal errorとなり、Goとの検証エラーになった。成功系の評価順序検証には、メソッド名を評価する関数内で例外を投げる形を使う。Goが直接`Class::{throw ...}()`を受理する差は、718件比較時点では残件だった。後続の上記727件比較で収録範囲の修正を確認した。[進捗記録](items6-10-progress.md)を参照。

## 項目7の宣言検査追補（2026-10-14、当時は追加fixtureが公式比較待ち）

具象クラスの登録時には、直接の`implements`だけでなく、親クラス連鎖で宣言されたinterfaceの契約も確認する。抽象親がメソッド実装を保留していても具象子はその契約を満たす必要がある。interface本体のメソッド検査は構文上の宣言順で行い、Goのmap反復順に診断対象を左右させない。trait合成時の`as`は元メソッドを残したまま別名を追加し、クラス本体の同名メソッドは優先する。異なるtraitメソッドの別名同士が同名なら競合として拒否する。既存のtrait定義を変更せず、利用側の宣言に複製する構造は維持した。根拠は[PHPのinterface規則](https://www.php.net/manual/en/language.oop5.interfaces.php)と[traitの優先順位・適応句](https://www.php.net/manual/en/language.oop5.traits.php)。未収録の宣言診断と複合型の共変・反変は残件。[実装・比較状況](items6-10-progress.md)を参照。

2026-10-14の再比較で、`set_include_path('')`は`false`を返し、現在の`include_path`を変更しないと判明した。変換後の空文字列を状態更新より先に判定する。[空文字fixtureの修正前公式ログ](validation/include-path-702-before-fix-mac.log)を根拠に期待値を訂正し、戻り値と状態をGo単体テストで検査した。その後の[修正後の公式PHP3版比較](validation/include-path-702-after-fix-mac.log)で当該fixtureを含む702件が各版で一致した。今回は保存済みログの確認であり、比較の再実行ではない。

2026-10-14追補：`set_include_path()`と`ini_restore()`の文字列引数は、評価時のファイルの`strict_types`を使って検査する。厳密モードのnullは変換・非推奨診断より先にTypeErrorとなり、`includePath`状態を変えない。弱い型モードでは従来どおり非推奨診断の後に空文字列へ変換する。根拠は[PHPの組み込み関数の型規則](https://www.php.net/manual/en/functions.internal.php)と[今回の回帰記録](items6-10-progress.md)。追補時には接続制限で実比較できなかったが、後続の[702件成功ログ](validation/include-path-702-after-fix-mac.log)で収録fixtureの一致を確認した。

## ローカルパス探索とinclude_path（2026-10-14）

`evaluator.includePath`を各`Run`で`.:/usr/local/lib/php`に初期化する。これは比較器が`php -n`で起動する固定Unixコンテナの初期値で、[684件の確定済み診断](validation/error-handler-684-fix-compare-all-mac.log)にも現れる。`get_include_path()`は現在値、`set_include_path()`は直前値を返す。PHP 8.0以降に`restore_include_path()`は存在しないので、復元には`ini_restore('include_path')`の最小実装を使う。他のphp.ini設定は管理しない。

`evalInclude`は式を一度評価し、`resolveLocalInclude`で4種の読み込み構文を共通に解決する。絶対パス、`./`、`../`はinclude_pathを使わず、通常の相対名はコロン区切りの各探索先を順に試した後に実行中ファイルのディレクトリを試す。相対探索先はプロセスの作業ディレクトリ基準。空のinclude_pathは直接の相対パスとして扱う。成功候補は実パスに変換するため、表記違いと相対シンボリックリンクが同じ`_once`履歴を参照する。履歴は既存どおり`*_once`を開いた直後、通常のinclude/requireを解析成功後に登録する。診断に示すinclude_pathは変更後の現在値を用いる。

探索順序は[PHPのinclude説明](https://www.php.net/manual/en/function.include.php)と[PHP 8.4の`php_resolve_path`](https://github.com/php/php-src/blob/PHP-8.4/main/fopen_wrappers.c#L475-L614)を根拠にした。`restore_include_path()`の削除は[公式マニュアル](https://www.php.net/manual/en/function.restore-include-path.php)、`ini_restore`の返り値と引数は[公式マニュアル](https://www.php.net/manual/en/function.ini-restore.php)に基づく。追加18件は[修正後の公式PHP3版比較](validation/include-path-702-after-fix-mac.log)で収録範囲の一致を確認した。空要素の探索・全権限エラー・ストリームラッパー・Windows形式のパスは後続で確認する。

## エラーハンドラと診断の経路（2026-10-14、収録fixtureは公式比較済み）

`evaluator`に現在のハンドラ、レベルマスク、登録履歴を保持し、呼び出し中のハンドラはその呼び出しのローカル状態へ一時退避する。これらは、独立した`Run`間では共有しない。診断名（Warning／Notice／Deprecated）からエラーレベルへ変換し、対象レベルだけを登録済みcallableへ渡す。ハンドラが`false`を返す場合は従来のstderr診断へ戻し、それ以外は通常診断を抑制する。実行中のハンドラは一時的に登録から外す。新規登録がなければ入れ子の警告は通常診断へ出すが、新規ハンドラが登録されればそのハンドラに渡す。呼び出し終了時にハンドラが未登録なら元のハンドラを戻す。実行中に`set_error_handler(null)`を呼んでも復元する。警告メッセージは切り詰めない。警告で投げられた例外は既存の関数・include・try/catchの`error`制御経路へ返し、未定義変数・配列キー・プロパティ読み取り、読み込み失敗と文字列変換の主要経路では後続評価を止める。まだすべての警告発生箇所で例外中断を検証していないため、この経路を完全対応とは扱わない。`trigger_error(E_USER_ERROR)`の捕捉されない場合は`fatalDiagnostic`でfinallyを通さず終了させ、ハンドラが処理した場合は実行を続ける。

追加22件を含む684件は、エラーハンドラ診断と未使用の単独変数式を修正した後の[Mac比較](validation/error-handler-684-fix-compare-all-mac.log)で公式PHP3版とも684一致・0不一致・0エラー、`make compare-all`終了0を確認した。修正前の不一致は[履歴](validation/error-handler-null-restore-compare-all-mac.log)に残す。[進捗](items6-10-progress.md)と追加前の[662件成功ログ](validation/throwable-chain-662-compare-all-mac.log)を参照。

## 例外連鎖とトレース（2026-10-14）

`exceptionData`は`message`・`code`・`file`・`line`、前の`Throwable`オブジェクト、生成時の`callFrame`スナップショットを保持する。`getPrevious()`はオブジェクトのポインタ同一性を保ち、`getTrace()`は同じフレームを配列へ写して返す。`getTraceAsString()`と未捕捉診断もそのフレームから整形し、表示文字列の逆解析はしない。比較器の`php -n`は`zend.exception_ignore_args`の既定値Offを用いるため、フレームの引数も保持する。配列値は不変の更新経路を共有し、オブジェクトはポインタを共有する。独立した`Run`ごとにフレームとクラス表を作る。

例外生成時に位置とトレースを固定し、別の場所でthrowまたは再throwしても書き換えない。組み込みの実行時エラーは捕捉可能な段階で同じ構造へ変換する。`finally`が先行例外を別の例外へ置き換えた場合は、新例外の既存`previous`鎖の末尾に先行例外をつなぐ。連鎖の追加・診断はオブジェクト同一性で循環を防ぐ。`finally`のreturnは通常例外を置き換え、捕捉できない宣言失敗は従来どおりfinallyへ流さない。[修正後の662件成功ログ](validation/throwable-chain-662-compare-all-mac.log)で公式PHP3版とも全一致し、641件以降の追加21件の比較待ちは解消した。今回は保存済みログの確認であり比較は再実行していない。[進捗](items6-10-progress.md)を参照。

PHP 8.4以降のクロージャのトレース表示名は、宣言元`sourceContext.filename`と`methodDecl.line`をフレーム生成時に組み込み、8.3では`{closure}`を維持する。同一フレームを配列形式と文字列形式に利用する。

トレースの各`callFrame`には評価済みの呼び出し引数と仮引数への束縛結果を別々に保持する。例外生成時、名前付き指定で省略された既定値を仮引数順に補い、可変長引数はローカル配列の要素に展開してフレームを固定する。これにより同じ構造化フレームから`getTrace()`と`getTraceAsString()`を生成する。派生例外の`parent::getMessage()`等は親クラス文脈の組み込みアクセサへ接続し、例外オブジェクトの同一性を維持する。この段階で追加した3 fixtureは後続の662件比較で3版とも一致した。

可変長引数については、ローカル配列が関数内で変更・追加・`unset`されても呼び出し時のトレースを保つため、型変換後の束縛値を`callFrame.boundArgs`へ固定する。スナップショット時にローカル配列を読み直さない。通常の仮引数と既定値の処理は従来どおりである。変更・追加・`unset`の3 fixtureはGoで検証済み、公式PHPとの比較待ちである。

## 読み込みと例外の統合（2026-10-14）

主ファイルは実行前に解析するため、その構文エラーは従来の解析診断として終了する。実行中に読み込んだファイルの構文エラーは、読み込み元の`try`で捕捉できる`ParseError`へ変換し、読み込み先のファイル名・行番号を`sourceError`で保持する。定数・関数等の宣言検査失敗は`fatalDiagnostic`として区別し、`catch`の対象にしない。`CompileError`を`Error`の子、`ParseError`を`CompileError`の子として組み込み階層に追加した。`include`／`require`の解決後は呼び出し元のファイル・行・演算子名を`callFrame`へ積み、解析・実行・例外伝播を通して保持し、必ず復元する。関数フレームと組み合わせて`getTraceAsString()`・未捕捉診断に利用する。

実パスへの正規化は既存の`loadedFiles`で継続する。ファイル未発見は履歴へ登録しない。`include_once`／`require_once`はファイルを開いた後、解析前に登録する。通常の`include`／`require`は解析・宣言検査に成功した後で登録するため、解析失敗後の`*_once`は再解析する。`*_once`での解析失敗後は同じ実パスへの`*_once`をスキップする。実行中の例外を捕捉した後も履歴は残る。読み込み先の`return`はinclude式の値へ戻し、関数の`return`や`finally`の制御とは区別する。実行中だけ`filename`と`strictTypes`を切り替えて戻し、変数スコープと`global`／`static`の束縛は読み込み元と共有する。致命的な宣言エラーは`Throwable`や通常のreturn制御へ変換せず、`finally`のreturnでも打ち消さない。現在の比較用公式PHP環境の既定include_pathを想定した失敗診断と、追加19件は後続の641件成功ログで公式PHP3版との一致を確認した。[進捗・再比較手順](items6-10-progress.md)を参照。

弱い型の`int`引数へ小数を渡して精度が失われた場合、呼び出し先の宣言行に非推奨警告を出す。PHP 8.5の読み込み先でのコンパイル時宣言エラーでは、捕捉したフレームから`require()`を除く。実行時エラーや例外のトレースには読み込みフレームを残す。

## 名前付き引数と例外の基盤（2026-10-14）

`namedArgument`は名前、式、行番号を保持する。`evalArguments`は明示引数と展開配列の挿入順に式を一度評価し、文字列キーを名前付き引数、整数キーを位置引数として`callArgument`へ記録する。`bindCallArguments`が仮引数名への束縛、重複・未知名、可変長引数のキー保持を担当する。関数・メソッド・staticメソッド・コンストラクタ・クロージャ・対応済み組み込み関数で同じ表現を使う。展開後の明示位置引数は解析時、展開内部で名前の後に整数キーが来る場合は実行時に拒否する。Traversableはまだ展開できない。束縛失敗の前後に評価される式の未収録の境界は未検証である。

公式PHPの修正前622件比較で、重複名・未知名の束縛エラーは呼び出し先のフレームを作らないと確認した。`invokeWithCapture`は`bindCallArguments`を呼び出し先スコープ・フレームの生成前に実行する。失敗時は呼び出し元フレームだけがトレースに残り、関数本体での失敗時は従来どおり呼び出し先フレームを残す。修正後は[ユーザーのMacでの全件比較](validation/named-exception-622-compare-all-mac.log)で各版622件一致した。今回はその保存済みログを確認したもので、比較を再実行していない。

`throwExpr`は評価値がThrowable階層のオブジェクトであることを確認し、`thrown`制御信号として伝播させる。`try`は通常のerror、`returned`、`loopControl`、`thrown`を区別し、対応するcatchがあればオブジェクトを局所スコープへ束縛する。finallyの制御信号があれば先行するreturn・break・continue・throwを上書きする。`runtimeError`のうち既存のTypeError等のクラスを持つものだけをcatch可能な例外オブジェクトへ変換し、構文・宣言診断は別経路に保つ。実行単位ごとにbuiltinクラスとハンドラのない状態を初期化する。未捕捉例外の全診断形式、Throwableの継承規則、警告ハンドラは未完成である。[収録fixtureと残件](items6-10-progress.md)を参照。

レビュー後、`try`節の有無は空ブロックのnilとは別に判定する。catch変数の書込先は直接のローカル写像ではなく`variableSlot`で解決し、宣言順によって選ばれたglobal／staticの束縛を保つ。名前付き引数後の引数展開は`parser.arguments`で式評価前に拒否し、未選択分岐にも適用する。

## union型の引数・戻り値（2026-10-15、収録fixtureは公式比較済み）

`typeDeclaration`は`|`で区切った構成型をソース順に名前解決し、重複と明らかな冗長型を宣言時に拒否する。`types.go`の`declaredType`は構成型のリストとnullableを分ける。実行時はまず値の型に一致する構成型を選び、弱い型モードで一致がなければ`int`→`float`→`string`→`bool`の変換候補を試す。数値文字列は整数形か小数形かで`int`と`float`の選択を変える。選択した型を既存の精度損失診断へ渡すので、警告ハンドラが例外を投げた時には関数本体に進まない。`invokeWithClass`のファイル単位strict_typesと、名前付き／可変長／展開済み引数の束縛はそのまま利用する。同じ構成型集合を持つ継承・interface署名は順序を問わず受理し、`?T`と`T|null`も同じ集合と扱う。unionを含む他の分散関係は現時点では受理しない。構成型を文字列からリストに分ける経路は将来のintersection／DNFで型項構造へ拡張する必要がある。[収録fixtureの比較結果と残件](items6-10-progress.md)を参照。

以下の節は各実装時点の設計履歴を残す。例外や名前付き引数、union型を「未実装」とする記述はその節の時点を指す。

## 単純型・ファイル単位strict_typesの基盤（2026-10-14）

`parameter`と`methodDecl`に単純型とnullableを保持し、`sourceContext`に宣言ファイルの`strict_types`を加えた。`parser.program`は先頭の`declare(strict_types=0/1)`を読み、主ファイル・読込先ごとに独立した設定にする。`include.go`は読込先の実行中だけ設定を切り替え、呼出し元へ戻る際に復元する。`invokeWithClass`は**引数には呼出し元ファイル**、**戻り値には宣言元ファイル**の設定を使う。`types.go`が単純型の照合・基本変換を担当し、失敗は既存の`runtimeError`に`TypeError`を付けて診断する。これはまだ`catch`可能なPHP例外オブジェクトではない。関数・メソッドの引数と戻り値だけが対象で、プロパティ／定数型、union・intersection・DNF、完全な型分散と弱い変換の全警告は残る。

interfaceの名前別メソッド表は探索用に保ち、宣言互換性検査では`interfaceContracts`で全祖先の宣言をたどる。これにより同名メソッドを持つ複数interfaceから一つだけを採用してしまう問題を防ぐ。文字列変数の`new`と`instanceof`はソース上の名前空間インポートを適用せず、値のクラス名を直接検索する。二要素配列callableは対象メソッドを一度解決してから既存の引数評価・呼出し経路へ渡す。値捕捉クロージャ・arrow関数・staticクロージャの部分集合を追加した。参照捕捉、一般のcallable判定、動的staticメンバーはまだない。

今回の[追加17件](validation/items6-10-types-new-fixtures.txt)はGoローカルテスト済みで、14件は公式PHP3版との比較で一致し、修正した3件は再比較待ち。[進捗と検証](items6-10-progress.md)を参照。既存588件の公式PHP全一致は追加前の履歴として維持する。クロージャ値は独立したオブジェクトに宣言・捕捉値・束縛された`$this`を持ち、呼び出し時に新しいフレームへ渡す。[PHP公式の匿名関数](https://www.php.net/manual/en/functions.anonymous.php)と[arrow関数](https://www.php.net/manual/en/functions.arrow.php)に従い、値捕捉は定義時の値を保存し、呼び出し中の変更を捕捉元へ書き戻さない。`static`ローカルは各クロージャ実体で保持する。可変長パラメータは末尾のみ受理し、余剰引数を配列へ束縛する。呼び出し側の`...`は数値キー配列を順に展開する。`parser.arguments`は展開済みの状態を構文解析中に保持し、その後の通常位置引数を実行前の致命エラーとして拒否する。複数の展開は許す。文字列キーは名前付き引数としての束縛が未実装なので明示的に拒否し、Traversableも未対応である。参照捕捉、完全な`Closure`組み込みAPI、引数の名前付き束縛は未実装である。

## 項目7のクラス登録とメンバー格納（2026-10-14）

`classDecl.kind`でclass/interface/traitを区別し、前方登録時に親だけでなくinterface・trait依存をたどる。実行途中のincludeや条件付き宣言のため、未知の依存は既存の保留方針に従う。interface契約は継承し、`instanceof`はクラスの親鎖とinterface鎖を調べる。型付き署名の検査は項目8の残件。

traitのメソッド・プロパティ・定数は利用先へコピーし、元のtraitを変更しない。コピーしたメソッドの`owner`は利用クラスで、privateプロパティの内部キーも利用クラスに結び付く。`insteadof`で元メソッドを選び、`as`で別名・可視性・finalを反映する。複数traitやクラス定数の全競合規則は未完成で、公式比較で検証が必要。

分割した`use`を処理する間、`methodDecl.traitSource`に取り込み元を残す。`decl.methods`に既に名前があっても、クラス本体の直接定義（`traitSource`が空）だけが後続traitの競合を解決する。先行するtrait由来なら競合として拒否する。これはメソッドの実行スコープを表す`owner`とは別の情報である。

クラスのstaticプロパティは`classDecl.staticProperties`のセルで保持する。継承では親セルを共有し、再宣言時は子のセルを作る。これは関数内の`static`ローカル変数の`staticScopes`とは独立する。クラス定数は値と宣言元・可視性・finalを保持する。呼出しフレームの`scope`はメソッド宣言元、`lateClass`は最後の非転送呼出しで決まるクラスとして分ける。`self::`・`parent::`・`static::`はメソッド探索先だけを変え、次の呼出しへ`lateClass`を転送する。明示的な`A::`は`lateClass`をAへ更新する。`new self`と`new static`も宣言元と呼出し元を区別する。traitの同名候補が複数あっても利用クラス自身の同名メソッドが優先される。独立Runは新しいクラス表とセルを作る。型情報・Reflection・autoload全体はまだ持たない。


## 過去の設計記録

以下は各段階で記録した当時の状態であり、「未実施」「未対応」はその時点を指す。現在の項目7の対応状況は上記と[588件の公式PHP全件一致記録](validation/item7-588-user-comparison-after-fix.log)を優先する。ただし、全一致は収録fixtureの範囲に限る。

同名privateプロパティ対応では、宣言した`propertyDecl.owner`を格納キーにも反映した。privateは`\0<宣言元クラス名>\0<プロパティ名>`、public/protectedは従来の単純名を`phpObject.properties`に用いる。親のメソッドで保持する`e.scope`が親なら親privateを選び、子メソッドなら子privateを選ぶ。子が親のpublic/protectedを合法に再宣言すると同じ単純名スロットを置き換える。アクセス可否と格納キーの決定を`propertyKey`へ集約し、読み取り・代入・quiet read・unsetを同じ規則へ通した。オブジェクト比較は値表の全スロットを見る。オブジェクト反復は現実装で未対応。private内部キーは利用者向け診断に出さない。新規fixtureの公式PHP比較はPodman接続拒否で未実施であり、細かな警告や診断の互換性は暫定判定である。

547件の公式PHP一致確認後、型なしクラス階層の宣言検査を部分追加した。親のpublic/protectedプロパティは子が同名で可視性を維持または拡大した場合に実効表の同じ位置で置き換える。メソッドも可視性を狭める上書きを拒否し、親のprivateメソッド内での同名呼出しは宣言元のprivate定義へ結びつける。クラスとメソッドにfinal・abstract情報を持たせ、継承・上書き、抽象メソッドの未実装、抽象クラス直接生成を検査する。親privateプロパティと子の同名宣言は独立した二つの格納場所が必要だが、現状の`phpObject.properties`は名前だけをキーにするため未対応として拒否する。クラス階層の型互換性、interface、trait、abstractの詳細制約、staticメンバーはまだ設計・実装の残件。新規3 fixtureはGo側のみ検証済みで、公式PHPの診断・時点は未確認である。

親クラス読込の再修正では、静的な`require`走査をクラス事前登録の根拠にする方法をやめた。関数呼び出しやswitchなどの実行文は、任意のコードを通じてクラス表を更新し得る境界とする。未知の親を持つ子がその境界より後にある場合は結合を保留し、宣言実行時に現在のクラス表から解決する。同一ファイルの無条件な前方宣言は依存順に事前登録する。この方針は既存の早期診断テストを維持するため、境界前の未知の親は早期に失敗する。境界を含む全ての診断時点は公式PHP比較後に再評価する。

2026-10-13のレビュー修正では、関数フレームの `staticNames` / `globalNames` を単一の `bindings` 表へ置き換えた。同名の宣言を実行した順に現在の格納先を切り替え、`unset` は現在の束縛を解除してルート・staticの値を保持する。参照一般の共有セルはまだない。クラス事前登録は同一ファイル内の無条件の親子を依存順に処理し、読込先や条件付き宣言で親が現れ得る子は宣言到達時に解決する。従来の前方参照を保ちつつ、`require` 後の継承を可能にする。588件の公式比較で、先行する実行文があるinterface実装・trait合成のエラーは宣言到達後に出ると確認されたため、その場合は事前登録を保留する。詳細な結果は[進捗](items6-10-progress.md)に記録した。

対象は PHP 8.4.26 を主基準とし、8.3.35 / 8.5.11 を同じ Go 実装の設定で扱う。PHP プログラムの評価を公式 PHP へ委譲しない。

変数は `variableScope` の名前から `value` への写像に置く。グローバルは `evaluator.variables`、通常関数とメソッドは呼び出しごとの `callFrame.locals` で、呼び出し終了時に呼び出し元を復元する。代入対象の `slot` はその写像を捕捉するので、評価途中で呼び出しが発生しても別フレームへ書き込まない。将来の参照はこの位置に共有セルを導入する余地があるが、現在の `value` は参照ではない。配列は更新経路をコピーする値、オブジェクトは共有するポインタとして維持する。

関数・クラス・定数は `runtimeSymbols` に登録する。関数とクラスのキーは大文字小文字を畳み、定数のキーは元の文字列を使う。ソース上の名前は構文解析時の `namespace` とカテゴリ別 `use` で解決し、修飾のない関数・定数のみ実行時にグローバル名へフォールバックする。文字列で渡された `define` / `defined` / `constant` の名前にはソースの `use` を適用しない。`sourceContext` は宣言した関数・クラス・メソッドにソースファイルと名前空間を保持する。現段階は単一ファイル実行のため診断のファイル名は CLI に渡されたファイル名である。

構文解析はファイル全体で実施し、未選択分岐も検査する。実行時の `runtimeError` は行・エラークラス・メッセージを持つ。`fatalDiagnostic` と構文エラーは別の経路で出力し、従来の呼び出し履歴を維持する。これは将来の例外処理の入力情報であり、`try` / `catch` や例外オブジェクト自体は未実装である。

2026-10-13の項目6〜10作業では、`evaluator.loadedFiles` を実行ごとに用意し、主ファイルと読み込み先の実パスを登録した。`*_once` はこの集合で判定し、循環中も再実行しない。読込元の変数スコープとシンボル表を共有し、読込先だけを別に構文解析・事前検査する。実行中は `filename` を切り替え、戻る際に復元する。関数呼び出しも宣言元ファイルを参照する。読込先の `return` は値を `include` 式へ返す。実行時エラーの出所は `sourceError` で保持する。深さ256の独自上限、include_path・ストリームラッパー・失敗診断の正確な互換性は残る。詳細は[項目6〜10の進捗](items6-10-progress.md)を参照。

`global` のため、呼び出しフレームはグローバル名の束縛表を持つ。通常の読み書きはローカル写像、`global` 宣言した名前の読み書きは実行単位のルート写像へ向ける。`static` ローカルは関数・メソッドの宣言ごとに保持する写像へ向け、再帰呼び出しでも共有する。`unset` は現在フレームの束縛だけを解除し、ルートまたは永続写像の値を残す。これは参照一般のセル基盤ではなく、参照捕捉・参照foreach等は未実装である。クラス階層、型メタデータ、クロージャ環境、例外制御フローもこの記録時点では項目7〜10の残件であり、実装したかのようには扱わない。

同じ作業で単純な親クラスのポインタを `classDecl` に追加し、親のメソッドとプロパティを子の実効表へ取り込む処理を試作した。メソッド・プロパティの宣言元を保持し、継承したメソッドの `$this` は子オブジェクトを指しつつ、可視性検査は宣言元のスコープで行う。親クラスへの `instanceof` はこの親ポインタをたどる。privateプロパティの同名再宣言やinterface・trait・final等は未対応であり、クラス階層全体の設計はまだ完了していない。型メタデータ、クロージャ環境、例外制御フローも項目8〜10の残件である。

既存の466件は着手前にソースと期待値を保存した。新機能追加後も全466件の保存済み期待値は変更していない。配列コピー、オブジェクト共有、評価・警告順、診断のファイル名・行番号を含む回帰確認は `make test` で行う。新規fixtureは公式 PHP 実比較が可能になるまで暫定期待値である。
## callable／Closureの基本解決（2026-10-16の実装、現在は収録範囲で公式比較済み）

`call_user_func_array`と`call_user_func`自身の引数も`evalArguments`の評価順を保ったまま`bindCallArguments`で仮引数名へ束縛する。これにより`args:`を`callback:`より先に指定しても、コールバックの解決と転送を束縛済みの値から行う。前者の転送配列では文字列キーを名前付き引数として保持し、後者の可変長引数も名前を保持する。この回帰2件を含む収録924件は[版別比較集計](validation/callable-924-after-fix-summary-mac.log)で3版とも一致した。修正直後の対象比較はPodman接続拒否で未実施だった。[進捗](items6-10-progress.md)。

`resolvedCallable`が、関数名文字列、`Class::method`文字列、二要素配列、Closureを現在のスコープから解決する。型検査、`is_callable`、`Closure::fromCallable`、エラーハンドラ登録、第一級callableの生成で同じ探索・可視性・static性の判定を使う。第一級callableは生成時に対象を確定し、`phpClosure`に元メソッド・受け手・called scopeを保持する。通常の無名関数／arrow関数は既存の宣言・値捕捉・`$this`を保持し、再呼び出しと配列・引数経由の同一性は`phpObject`を共有する。`call_user_func`系は評価済み引数を既存の束縛へ渡す。組み込み関数へ転送する場合は既存の組み込み引数処理を使う。`callable`は引数・戻り値に許可し、状態により呼び出し可能性が変わるプロパティ・昇格プロパティでは拒否する。詳細と公式比較の状態は[進捗](items6-10-progress.md)。

`Closure::bind`／`bindTo`／`call`、参照捕捉、`__invoke`、マジック呼び出し、組み込みcallbackの全面対応は未実装。追加16件を含む収録924件は版別集計で一致したが、これら未収録の機能・境界条件は引き続き残る。
