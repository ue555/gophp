<?php
class QuietBox { public ?int $value; }
$box = new QuietBox;
$name = 'value';
echo isset($box->$name) ? 'yes' : 'no';
echo '|', empty($box->{$name}) ? 'empty' : 'full';
echo '|', $box->{$name} ?? 'missing';
echo '|', ($box->{$name} ??= '6');
echo '|', $box->value;
unset($box->{$name});
echo '|', isset($box->{$name}) ? 'yes' : 'no';
echo '|', $box->{$name} ?? 'missing', "\n";
