#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
IPv6地址简写转换脚本
支持批量处理CSV、XLS、XLSX文件中的IPv6地址，将完整格式转换为RFC 5952规范简写格式
"""

import pandas as pd
import ipaddress
import re
import sys
import argparse
import io
from pathlib import Path
from typing import Optional
from tqdm import tqdm

# Windows 控制台编码处理
if sys.platform == 'win32':
    try:
        import locale
        # 设置标准输出为 UTF-8
        sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8', errors='replace')
        sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8', errors='replace')
    except:
        pass
    try:
        # 设置 Windows 控制台代码页为 UTF-8
        import ctypes
        kernel32 = ctypes.windll.kernel32
        kernel32.SetConsoleMode(kernel32.GetStdHandle(-11), 7)
        kernel32.SetConsoleOutputCP(65001)
        kernel32.SetConsoleCP(65001)
    except:
        pass


# IPv6正则表达式（用于快速识别）
IPV6_PATTERN = re.compile(
    r'(?:(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|'  # 完整格式
    r'(?:[0-9a-fA-F]{1,4}:){1,7}:|'                    # 以::结尾
    r':(?:(?:[0-9a-fA-F]{1,4}:){1,7}|)|'              # 以::开头
    r'[0-9a-fA-F]{1,4}::(?:[0-9a-fA-F]{1,4}:){0,5}[0-9a-fA-F]{1,4}|'  # 中间::
    r'::(?:[0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4})'  # ::在开头或中间
)


def is_ipv6(value: str) -> bool:
    """快速判断字符串是否为IPv6地址"""
    if not isinstance(value, str) or not value.strip():
        return False
    return bool(IPV6_PATTERN.match(value.strip()))


def compress_ipv6(ip_str: str) -> Optional[str]:
    """
    将IPv6地址转换为简写格式（RFC 5952规范）
    - 移除前导零
    - 使用::压缩最长的连续零块
    """
    if not isinstance(ip_str, str):
        return ip_str

    ip_str = ip_str.strip()
    if not is_ipv6(ip_str):
        return ip_str

    try:
        # 使用ipaddress模块自动处理简写（符合RFC 5952）
        ip_obj = ipaddress.IPv6Address(ip_str)
        # compressed格式会自动进行最优压缩
        return ip_obj.compressed
    except (ipaddress.AddressValueError, ValueError):
        return ip_str


def find_ipv6_columns(df: pd.DataFrame) -> list:
    """自动识别包含IPv6地址的列"""
    ipv6_cols = []
    sample_size = min(1000, len(df))  # 采样检查，提高性能

    for col in df.columns:
        # 采样检查该列是否包含IPv6
        sample = df[col].dropna().head(sample_size)
        if sample.empty:
            continue

        ipv6_count = sum(1 for val in sample if is_ipv6(str(val)))
        ipv6_ratio = ipv6_count / len(sample)

        # 如果超过30%的值是IPv6，认为是IPv6列
        if ipv6_ratio > 0.3:
            ipv6_cols.append(col)

    return ipv6_cols


def process_file(
    input_path: str,
    output_path: Optional[str] = None,
    columns: Optional[list] = None,
    column_pattern: Optional[str] = None,
    show_progress: bool = True
) -> pd.DataFrame:
    """
    处理文件中的IPv6地址

    Args:
        input_path: 输入文件路径
        output_path: 输出文件路径（可选，默认在原文件名加_compressed）
        columns: 指定要处理的列名列表（可选）
        column_pattern: 列名匹配模式（正则，可选）
        show_progress: 是否显示进度条

    Returns:
        处理后的DataFrame
    """
    input_path = Path(input_path)

    if not input_path.exists():
        raise FileNotFoundError(f"文件不存在: {input_path}")

    # 确定输出路径
    if output_path is None:
        suffix = input_path.suffix
        output_path = input_path.with_stem(f"{input_path.stem}_compressed").with_suffix(suffix)
    else:
        output_path = Path(output_path)

    print(f"正在读取文件: {input_path}")

    # 根据文件类型读取
    if input_path.suffix.lower() == '.csv':
        df = pd.read_csv(input_path, low_memory=False)
    else:
        # xlsx/xls文件
        df = pd.read_excel(input_path, engine='openpyxl' if input_path.suffix.lower() == '.xlsx' else 'xlrd')

    print(f"文件读取完成，共 {len(df)} 行, {len(df.columns)} 列")

    # 确定要处理的列
    if columns:
        target_cols = columns
        print(f"使用指定列: {target_cols}")
    elif column_pattern:
        pattern = re.compile(column_pattern)
        target_cols = [col for col in df.columns if pattern.search(str(col))]
        print(f"匹配到列: {target_cols}")
    else:
        target_cols = find_ipv6_columns(df)
        if target_cols:
            print(f"自动识别到IPv6列: {target_cols}")
        else:
            print("警告: 未自动识别到IPv6列，将尝试处理所有文本列")
            target_cols = df.select_dtypes(include=['object']).columns.tolist()

    if not target_cols:
        print("警告: 没有找到需要处理的列")
        return df

    # 处理每一列
    for col in target_cols:
        if col not in df.columns:
            print(f"警告: 列 '{col}' 不存在，跳过")
            continue

        print(f"\n正在处理列: {col}")

        # 使用向量化操作优化性能
        # 先转换为字符串类型
        str_series = df[col].astype(str)

        # 使用progress bar
        if show_progress:
            tqdm.pandas(desc=f"  压缩 {col}")
            df[col] = str_series.progress_apply(compress_ipv6)
        else:
            df[col] = str_series.apply(compress_ipv6)

    # 保存结果
    print(f"\n正在保存到: {output_path}")
    if output_path.suffix.lower() == '.csv':
        df.to_csv(output_path, index=False)
    else:
        df.to_excel(output_path, index=False, engine='openpyxl')

    print(f"处理完成! 输出文件: {output_path}")
    return df


def main():
    parser = argparse.ArgumentParser(
        description='批量转换Excel/CSV中的IPv6地址为简写格式',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 自动处理
  python ipv6_compressor.py data.xlsx

  # 指定输出文件
  python ipv6_compressor.py data.xlsx -o result.xlsx

  # 指定列名
  python ipv6_compressor.py data.xlsx -c ip_address source_ip

  # 使用列名模式匹配
  python ipv6_compressor.py data.xlsx -p "ip|address"

  # 处理CSV文件
  python ipv6_compressor.py data.csv
        """
    )

    parser.add_argument('input', help='输入文件路径 (CSV/XLS/XLSX)')
    parser.add_argument('-o', '--output', help='输出文件路径')
    parser.add_argument('-c', '--columns', nargs='+', help='指定要处理的列名')
    parser.add_argument('-p', '--pattern', help='列名匹配模式（正则表达式）')
    parser.add_argument('--no-progress', action='store_true', help='不显示进度条')

    args = parser.parse_args()

    try:
        process_file(
            input_path=args.input,
            output_path=args.output,
            columns=args.columns,
            column_pattern=args.pattern,
            show_progress=not args.no_progress
        )
    except Exception as e:
        print(f"错误: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()
