<?php
function state($v) { echo json_encode([$v,preg_last_error(),preg_last_error_msg()]),"\n"; }
state(null);
state(@preg_match('/(/', 'x'));
state(preg_match('/x/', 'x'));
state(preg_match('//u', "\xff"));
state(preg_match('/./u', '日本', $m, 0, 1));
state(preg_match('/z/', 'x'));
state(@preg_replace('/(/','x','a'));
state(preg_replace('/a/','b','a'));
state(@preg_split('/(/','a'));
state(preg_split('/,/', 'a,b'));
state(preg_match_all('/./u',"\xff",$m));
state(preg_match_all('/./u','日本',$m));
state(@preg_replace_callback('/(/', fn($m)=>'z','x'));
state(preg_replace_callback('/x/', fn($m)=>'z','x'));
echo json_encode([PCRE_VERSION_MAJOR >= 10, PCRE_VERSION_MINOR >= 0, is_string(PCRE_VERSION), is_bool(PCRE_JIT_SUPPORT)]),"\n";
