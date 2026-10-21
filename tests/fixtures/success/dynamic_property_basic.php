<?php
class DynamicBox { public int $value = 1; public array $items = [4]; }
$box = new DynamicBox;
$name = 'value';
echo $box->$name, '|', $box->{$name}, '|';
echo ($box->{$name} = '2'), '|', ($box->{$name} .= '3'), '|';
echo ++$box->{$name}, '|', $box->{$name}++, '|', $box->value, '|';
$items = 'items';
$box->{$items}[] = 8;
$box->{$items}[0] = 5;
echo $box->{$items}[0], ',', $box->{$items}[1], "\n";
