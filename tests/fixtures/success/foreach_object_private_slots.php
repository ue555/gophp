<?php
class IteratedParent {
    private $same = 'parent';
    public function parentView() {
        foreach ($this as $key => $value) echo $key, ':', $value, ';';
    }
}
class IteratedChild extends IteratedParent {
    private $same = 'child';
    public function childView() {
        foreach ($this as $key => $value) echo $key, ':', $value, ';';
    }
}
$item = new IteratedChild;
foreach ($item as $key => $value) echo $key, ':', $value, ';';
echo '|';
$item->parentView();
echo '|';
$item->childView();
