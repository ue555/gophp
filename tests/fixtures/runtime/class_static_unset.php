<?php
class StaticUnset { public static $value = 1; }
echo "before\n";
if (false) { unset(StaticUnset::$value); }
echo "after\n";
