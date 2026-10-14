<?php
class OffsetReceiver { function take($value) { return "value" . $value; } }
class OffsetBox { public $items; }
$o = new OffsetBox;
$o->items = [new OffsetReceiver];
$i = 0;
$j = 0;
echo $o?->items[$i++]?->take($j++) ?? "fallback";
echo "|", $i, "|", $j, "\n";
