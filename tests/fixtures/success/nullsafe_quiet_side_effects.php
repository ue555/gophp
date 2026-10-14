<?php
class Source {
    public $items = ["A", "B", "C"];
    function take($value) { return [$value]; }
}
$o = new Source;
$null = null;
$i = 0;
echo $o?->items[$i++] ?? "fallback";
echo "|", isset($o?->items[$i++]) ? "yes" : "no";
echo "|", empty($o?->items[$i++]) ? "empty" : "full";
echo "|", $i;
echo "|", $null?->take($i++)[0] ?? "fallback";
echo "|", $i;
echo "|", $o?->take($i++)[0];
echo "|", $i, "\n";
