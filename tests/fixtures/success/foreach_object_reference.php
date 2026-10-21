<?php
class IteratedReference {
    public int $number = 1;
    public $text = 'a';
}
$item = new IteratedReference;
foreach ($item as $key => &$value) {
    if ($key === 'number') $value = 3;
    else $value = 'b';
}
echo $item->number, '|', $item->text, '|';
$value = 'c';
echo $item->text;
unset($value);
