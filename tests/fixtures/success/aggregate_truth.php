<?php
$o = new stdClass;
echo (![]) + 0, (!![false]) + 0, (!!$o) + 0;
echo ([] == false) + 0, ([] == null) + 0, ([0] == true) + 0, ($o == true) + 0, ($o != null) + 0;
echo ([] || [0]) + 0, ($o && [1]) + 0, "\n";
if ([]) echo 'bad'; elseif ([0]) echo 'array';
if ($o) echo '-object'; else echo 'bad';
echo "\n";
