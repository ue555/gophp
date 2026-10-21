<?php
class ChainDynamicBox { public array $items = [6]; }
$calls = 0;
function chainName() { global $calls; $calls++; return 'items'; }
function chainIndex() { global $calls; $calls++; return 0; }
$box = null;
echo $box?->{chainName()}[chainIndex()] ?? 'fallback';
echo '|', $calls, '|';
$box = new ChainDynamicBox;
echo $box?->{chainName()}[chainIndex()], '|', $calls, "\n";
