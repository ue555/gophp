<?php
class ParentResult { public function result(): mixed { return null; } }
echo "before\n";
class ChildResult extends ParentResult { public function result(): void {} }
echo "after\n";
