<?php
class IteratedBase {
    public $open = 'a';
    protected $family = 'b';
    private $hidden = 'c';
    public ?int $pending;

    public function inside() {
        foreach ($this as $key => $value) echo $key, ':', $value, ';';
    }
}
$item = new IteratedBase;
foreach ($item as $key => $value) echo $key, ':', $value, ';';
echo '|';
$item->inside();
