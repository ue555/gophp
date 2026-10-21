<?php
interface AnonymousContract { public function value(): int; }
trait AnonymousTrait { public function value(): int { return $this->number; } }
function makeAnonymous($number) {
    return new class($number) implements AnonymousContract {
        use AnonymousTrait;
        public function __construct(public int $number) {}
    };
}
$first = makeAnonymous(4);
$second = makeAnonymous(8);
echo $first->value(), '|', $second->value(), '|', ($first instanceof $second ? 'same' : 'other');
