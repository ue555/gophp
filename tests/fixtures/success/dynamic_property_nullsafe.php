<?php
class SafeBox { public int $value = 8; }
$calls = 0;
function safeName() { global $calls; $calls++; return 'value'; }
$box = null;
echo $box?->{safeName()} ?? 'fallback';
echo '|', isset($box?->{safeName()}) ? 'yes' : 'no';
echo '|', $calls, '|';
$box = new SafeBox;
echo $box?->{safeName()}, '|', $calls, "\n";
