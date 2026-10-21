<?php
function tokens($text,$flags=0){$out=[];foreach(token_get_all($text,$flags) as $t){$out[]=is_array($t)?[token_name($t[0]),$t[1],$t[2]]:$t;}echo json_encode($out,JSON_UNESCAPED_SLASHES),"\n";}
tokens('<?php public(set) string $p; $x |> trim(...); (void) $x; __PROPERTY__;');
tokens('<?php class C { public function match() {} const MATCH = 1; } C::match();');
tokens('<?php class C { public function match() {} const MATCH = 1; } C::match();',TOKEN_PARSE);
tokens("<?php\n\$x=<<<TXT\n    abc \$name\n    TXT;\n");
tokens('<?php $x = "{$a["key"]}";');
tokens('<?php /* a */ & /* comment */ $x; b"raw"; echo 1;?>tail');
