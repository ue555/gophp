<?php
interface A {}
interface B {}
interface C {}
echo "before|";
function redundantReordered((B&A)|(C&A&B) $value) {}
echo "after";
