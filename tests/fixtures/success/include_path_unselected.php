<?php
$calls = 0;
function sideEffectPath() { global $calls; $calls++; return 'missing-path.inc'; }
if (false) { include sideEffectPath(); set_include_path('invalid'); }
echo $calls, '|', (get_include_path() === '.:/usr/local/lib/php' ? 'default' : 'changed'), "\n";
