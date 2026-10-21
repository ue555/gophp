<?php
class ChangedReceiverBox { public int $value; }
$first = new ChangedReceiverBox; $first->value = 1;
$second = new ChangedReceiverBox; $second->value = 2;
$selected = $first;
function changeMember() { global $selected, $second; $selected = $second; return 'value'; }
$selected->{changeMember()} = 9;
echo $first->value, '|', $second->value, '|';
$selected = $first;
function changeRight() { global $selected, $second; $selected = $second; return 10; }
$name = 'value';
$selected->{$name} = changeRight();
echo $first->value, '|', $second->value, "\n";
