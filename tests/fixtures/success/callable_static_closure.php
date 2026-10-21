<?php
class ClosureHost {
	public string $label = 'bound';
    public function make(): array {
        $normal = function (): string { return $this->label; };
        $static = static function (): string { return 'static'; };
        return [$normal, $static];
    }
}
$parts = (new ClosureHost)->make();
echo $parts[0](), '|', $parts[1]();
