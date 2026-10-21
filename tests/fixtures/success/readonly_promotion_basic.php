<?php
class PromotedIdentity {
    public function __construct(public readonly int $id, protected string $label = 'default', private ?string $note = null) {
        echo $this->id, ':', $this->label, ':', $this->note ?? 'none', '|';
    }
    public function details() { return $this->label . ':' . ($this->note ?? 'none'); }
}
$o = new PromotedIdentity(note: 'memo', id: 7);
echo $o->id, ':', $o->details(), "\n";
