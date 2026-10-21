<?php
class ScopeBase {
    public function own(): self { return new self; }
    public function late(): static { return new static; }
}
class ScopeChild extends ScopeBase {}
echo (new ScopeChild)->own() instanceof ScopeChild ? 'wrong' : 'base';
echo '|';
echo (new ScopeChild)->late() instanceof ScopeChild ? 'child' : 'wrong';
