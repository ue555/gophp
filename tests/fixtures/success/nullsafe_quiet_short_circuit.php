<?php
$x = null;
$i = 0;
echo $x?->items[$i++] ?? "fallback";
echo "|", $i;
echo "|", isset($x?->items[$i++]) ? "yes" : "no";
echo "|", empty($x?->items[$i++]) ? "empty" : "full";
echo "|", $i, "\n";
