<?php
$o = null;
$i = 0;
echo $o?->secret?->take($i++) ?? "fallback";
echo "|", $i, "\n";
