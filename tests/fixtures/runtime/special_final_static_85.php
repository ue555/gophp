<?php
class FinalStaticContract { public function make(): static { return new static; } }
final class FinalStaticImplementation extends FinalStaticContract {
    public function make(): self { return new self; }
}
echo (new FinalStaticImplementation)->make() instanceof FinalStaticImplementation ? 'final' : 'wrong';
