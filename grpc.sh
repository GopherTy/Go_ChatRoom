#!/bin/bash
#Program:
#       自動 查找 編譯 google's grpc 代碼
#History:
#       2018-09-18 king first release support golang/c++
#       2019-02-24 king support java
#       2019-02-28 king support dart
#Email:
#       zuiwuchang@gmail.com

# 顯示幫助信息
function ShowHelp(){
    echo "help              : show help"
    echo "go root output    : build src/*.proto for golang" 
    echo "dart root output  : build src/*.proto for dartlang"
    echo "cpp root output   : build src/*.proto for c++"
    echo "java root output  : build src/*.proto for java"
}

# 遞歸 查詢所有的 檔案夾/包
# $1 根本目錄
function find_package(){
    _DIRS=""
    files=`find $1 -type d`
    ok=$?
    if [[ $ok != 0 ]];then
        exit $ok
    fi

    for str in $files
    do
        str=${str#$1}
        str=${str#.}
        str=${str#/}
        if [ "$str" ];then
            _DIRS="$_DIRS $str"
        fi
    done
}
# 查找 proto 檔案
# $1 根本目錄
# $2 包路徑
function find_proto(){
    _SOURCES=""

    files=`find $1/$2 -maxdepth 1 -name *.proto -type f`
    ok=$?
    if [[ $ok != 0 ]];then
        exit $ok
    fi

    for str in $files
    do
        str=${str#$1}
        str=${str#.}
        str=${str#/}
        if [ "$str" ];then
            _SOURCES="$_SOURCES $str"
        fi
    done
}

function print_source(){
    if [ "$2" ];then
        for str in $2
        do
            str=${str#$1}
            str=${str#.}
            str=${str#/}
            echo "   $str"
        done
    else
        echo "  warning : not found any source"
    fi
    echo "}"
    echo
}
function check_params(){
    if [ ! "$2" ];then
        echo "need param root directory"
        echo "exmaple : grpc.sh $1 proto protocol/$1"
        exit 1
    fi
    if [ ! "$3" ];then
        echo "need param output"
        echo "exmaple : grpc.sh $1 proto protocol/$1"
        exit 1
    fi

    if [ ! -d "$3" ];then
        echo "directory not exist : $3"
        exit 1
    fi
}
# 自動 查找 並編譯 grpc 到 go 代碼
# $1 protoc 根目錄
# $2 輸出目錄
function BuildGo(){
    check_params go $1 $2
    root=$1
    out=$2

    find_package $root
    for dir in $_DIRS
    do
        echo
        echo "package $dir {"
        find_proto $root $dir
        print_source $dir $_SOURCES
        if [ "$_SOURCES" ];then
            echo "protoc -I $root --go_out=plugins=grpc:$out $_SOURCES"
            protoc -I $root --go_out=plugins=grpc:$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi
        fi
    done
}

# 自動 查找 並編譯 grpc 到 dart 代碼
# $1 protoc 根目錄
# $2 輸出目錄
function BuildDart(){
    check_params dart $1 $2
    root=$1
    out=$2

    find_package $root
    for dir in $_DIRS
    do
        echo
        echo "package $dir {"
        find_proto $root $dir
        print_source $dir $_SOURCES
        if [ "$_SOURCES" ];then
            echo "protoc -I $root --dart_out=grpc:$out $_SOURCES"
            protoc -I $root --dart_out=grpc:$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi
        fi
    done
}

# 自動 查找 並編譯 grpc 到 c++ 代碼
# $1 protoc 根目錄
# $2 輸出目錄
function BuildCpp(){
    check_params cpp $1 $2
    root=$1
    out=$2

    GrpcCppPlugin=`which grpc_cpp_plugin`
    if [ -f "$GrpcCppPlugin".exe ];then
        GrpcCppPlugin="$GrpcCppPlugin".exe
    fi

    if [ ! -f "$GrpcCppPlugin" ];then
        echo grpc_cpp_plugin not found
        exit 1
    fi

    find_package $root
    for dir in $_DIRS
    do
        echo
        echo "package $dir {"
        find_proto $root $dir
        print_source $dir $_SOURCES
        if [ "$_SOURCES" ];then
            # pb
            echo "protoc -I $root --cpp_out=$out $_SOURCES"
            protoc -I $root --cpp_out=$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi

            # grpc
            echo "protoc -I $root --plugin=protoc-gen-grpc=$GrpcCppPlugin --grpc_out=$out $_SOURCES"
            protoc -I $root --plugin=protoc-gen-grpc=$GrpcCppPlugin --grpc_out=$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi
        fi
    done
}
# 自動 查找 並編譯 grpc 到 java 代碼
# $1 protoc 根目錄
# $2 輸出目錄
function BuildJava(){
    check_params java $1 $2
    root=$1
    out=$2

    find_package $root
    for dir in $_DIRS
    do
        echo
        echo "package $dir {"
        find_proto $root $dir
        print_source $dir $_SOURCES
        if [ "$_SOURCES" ];then
            # pb
            echo "protoc -I $root --java_out=$out $_SOURCES"
            protoc -I $root --java_out=$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi

            # grpc
            echo "protoc --plugin=protoc-gen-grpc-java --proto_path=$root --grpc-java_out=$out $_SOURCES"
            protoc --plugin=protoc-gen-grpc-java --proto_path=$root --grpc-java_out=$out $_SOURCES
            ok=$?
            if [[ $ok != 0 ]];then
                exit $ok
            fi
        fi
    done
}

ok=0
case $1 in
    go)
        BuildGo $2 $3
        ok=$?
    ;;

    dart)
        BuildDart $2 $3
        ok=$?
    ;;

    cpp)
        BuildCpp $2 $3
        ok=$?
    ;;

    java)
        BuildJava $2 $3
        ok=$?
    ;;

    *)
        ShowHelp
        ok=$?
    ;;
esac
exit $ok
