<?php
$o = new stdClass;
$o->items = [0 => 1];
echo isset($o?->items[$undefined + 0]) ? "yes" : "no";
echo "|", $o?->items[$alsoUndefined + 0] ?? "fallback", "\n";
