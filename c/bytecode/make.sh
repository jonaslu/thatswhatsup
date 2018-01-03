case $1 in
  "")
    cd vm && python setup_bytecode_ext.py build_ext --inplace
    ;;
  "clean")
    (cd vm; rm *.so; rm bytecode_ext.c; rm -rf build/)
  ;;
esac
