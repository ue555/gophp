<?php
class Reader { function take($value) { return $value; } }
$o = new Reader;
$null = null;
echo $o?->take($undefined + 0) ?? "fallback";
echo "|", $null?->take($skipped + 0) ?? "fallback", "\n";
