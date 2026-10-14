<?php
$o = new stdClass;
$o->items = ["keep" => 1, "drop" => 2];
unset($o->items["drop"]);
echo isset($o?->items["drop"]) ? "yes" : "no";
echo "|", $o?->items["keep"] ?? "missing", "\n";
