<?php
trait PromotedLabel {
    public function __construct(protected readonly string $label) {}
    public function label() { return $this->label; }
}
class ParentPromoted {
    use PromotedLabel;
}
class ChildPromoted extends ParentPromoted {}
$c = new ChildPromoted('hello');
echo $c->label(), "\n";
