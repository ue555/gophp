<?php
interface IntegratedLabel { public function label(): string; }
trait IntegratedLabelTrait {
    public function label(): string { return static::class; }
}
class IntegratedBase {
    public function inherited(): self { return new self; }
}
class IntegratedChild extends IntegratedBase implements IntegratedLabel {
    use IntegratedLabelTrait;
    private function secret(int $n): int { return $n + 3; }
    public function callback(): Closure { return $this->secret(...); }
}
function integratedNumber(int $n): int { return $n + 1; }
return 'loaded';
