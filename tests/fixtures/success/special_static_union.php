<?php
class UnionScopeBase { public function make(): static|null { return new static; } }
class UnionScopeChild extends UnionScopeBase {}
echo (new UnionScopeChild)->make() instanceof UnionScopeChild ? 'union' : 'wrong';
