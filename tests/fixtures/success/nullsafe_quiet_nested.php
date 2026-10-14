<?php
$o = new stdClass;
$o->items = ["child" => ["value" => 7]];
$i = 0;
$k = "child";
echo $o?->items[$k]["value"] ?? "fallback";
echo "|", isset($o?->items[$k]["value"]) ? "yes" : "no";
echo "|", empty($o?->items[$k]["value"]) ? "empty" : "full";
echo "|", $o?->items[$k]["missing"] ?? "fallback";
echo "|", isset($o?->missing->field) ? "yes" : "no";
echo "|", empty($o?->missing->field) ? "empty" : "full";
echo "|", $o?->missing->field ?? "fallback";
echo "|", $o?->missing[$i++] ?? "fallback";
echo "|", $i, "\n";
