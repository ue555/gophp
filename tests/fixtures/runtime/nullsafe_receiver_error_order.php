<?php
class LockedBox { private $secret = null; }
function argument() { echo "argument"; return 1; }
$o = new LockedBox;
echo "before|";
echo $o?->secret?->take(argument()) ?? "fallback";
