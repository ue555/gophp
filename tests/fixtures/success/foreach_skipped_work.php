<?php
false && isset($a[missing()]);
if (false) foreach (missing() as $v) echo 1 / 0;
foreach ([1] as $v) { break; echo missing(); }
foreach ([] as $v) { $bad = 1 / 0; }
echo $v, "\n";
