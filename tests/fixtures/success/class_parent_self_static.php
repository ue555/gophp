<?php
class ParentFactory {
    public const VALUE = "parent";
    public function value() { return self::VALUE; }
    public function makeSelf() { return new self; }
    public function makeLate() { return new static; }
}
class ChildFactory extends ParentFactory {
    public const VALUE = "child";
    public function value() { return parent::value() . "/" . static::VALUE; }
}
$o = new ChildFactory;
echo $o->value(), '|', ($o->makeSelf() instanceof ParentFactory), '|';
echo ($o->makeSelf() instanceof ChildFactory), '|', ($o->makeLate() instanceof ChildFactory), "\n";
