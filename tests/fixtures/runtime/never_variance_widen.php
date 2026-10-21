<?php
class ParentStop { public function stop(): never { throw new Exception('stop'); } }
echo "before\n";
class ChildStop extends ParentStop { public function stop(): void {} }
echo "after\n";
