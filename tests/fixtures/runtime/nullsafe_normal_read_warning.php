<?php
$o = new stdClass;
echo $o?->missing ?? "fallback";
echo "|";
echo $o?->missing;
echo "after\n";
