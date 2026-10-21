<?php
class UnselectedDynamicBox { public int $value = 3; }
$box = new UnselectedDynamicBox;
function unusedName() { echo 'name'; return 'value'; }
if (false) { $box->{unusedName()} = 9; }
echo $box->value, "\n";
