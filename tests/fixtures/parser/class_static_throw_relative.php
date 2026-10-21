<?php
class ParentTarget { public static function run() {} }
class ChildTarget extends ParentTarget {
    public static function neverCalled() {
        parent::{throw new Exception('stop')}();
    }
}
