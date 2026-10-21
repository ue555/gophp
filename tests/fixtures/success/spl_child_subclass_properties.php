<?php
class MarkedDirectory extends RecursiveDirectoryIterator {
    private int $marker = 7;
    public function marker(): int { return $this->marker; }
}
$root = new MarkedDirectory('tests/fixtures/support/spl_walk', RecursiveDirectoryIterator::SKIP_DOTS);
for ($root->rewind(); $root->valid(); $root->next()) {
    if ($root->hasChildren()) {
        $child = $root->getChildren();
        echo $child->marker(), '|', $child->valid() ? 'valid' : 'empty';
        break;
    }
}
