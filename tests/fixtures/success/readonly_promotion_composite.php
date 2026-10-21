<?php
interface PromotedA {}
interface PromotedB {}
class PromotedBoth implements PromotedA, PromotedB {}
class PromotedComposite {
    public function __construct(public readonly (PromotedA&PromotedB)|null $item, public int|string $id) {}
}
$a = new PromotedComposite(new PromotedBoth, 4);
$b = new PromotedComposite(null, 'next');
echo gettype($a->item), ':', $a->id, '|', gettype($b->item), ':', $b->id, "\n";
