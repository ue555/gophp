<?php
function dumpTokens($code, $flags=0) {
 $out=[]; foreach (token_get_all($code,$flags) as $t) { $out[]=is_array($t)?[token_name($t[0]),$t[1],$t[2]]:$t; }
 echo json_encode($out,JSON_UNESCAPED_UNICODE|JSON_UNESCAPED_SLASHES),"\n";
}
dumpTokens("before\r\n<?php \n/** doc */ // line\n# last\n\$x=1.2+0x10; echo '日'; ?>\r\nafter<?= 3 ?>");
dumpTokens('<?php use A\B; \A\B::class; namespace\x(); $a?->match($b & $c, ...$args); yield from $it; (int) $x;');
dumpTokens('<?php $x = "Hello $name {$user->name} ${old} $a[0]";');
dumpTokens("<?php\n\$x=<<<TXT\nhello \$name\nTXT;\n\$y=<<<'END'\nraw \$name\nEND;\n");
dumpTokens('<?php echo 1; ?>text', TOKEN_PARSE);
try {token_get_all('<?php echo ;',TOKEN_PARSE);} catch(ParseError $e) {echo get_class($e),': ',$e->getMessage(),"\n";}
dumpTokens('<?php echo ;');
echo json_encode([token_name(0),token_name(T_DOUBLE_COLON),T_INLINE_HTML,T_OPEN_TAG,T_WHITESPACE]),"\n";
