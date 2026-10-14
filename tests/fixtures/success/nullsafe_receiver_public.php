<?php
class PublicReceiver { function take($value) { return "value" . $value; } }
class PublicBox { public $child = null; }
$o = new PublicBox;
$i = 0;
echo $o?->child?->take($i++) ?? "fallback";
echo "|", $i, "|";
$o->child = new PublicReceiver;
echo $o?->child?->take($i++), "|", $i, "\n";
