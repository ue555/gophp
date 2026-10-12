<?php
$a = null;
echo isset($a[($a = [1])[0] - 1]) + 0, ':', $a[0], "\n";
