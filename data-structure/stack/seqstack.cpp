#include <iostream>

#define InitSize 10
#define IncSize 5

template <typename T>
class SeqStack {
public:
    SeqStack(int length = InitSize);  // 构造函数
    ~SeqStack();                      // 析构函数

public:
    bool Push(const T& e);  // 入栈
    bool Pop(T &e);  // 出栈 删除栈顶数据
    bool GetTop(T &e); // 读取栈顶元素

    void DispList();  // 输出栈中所有元素
    int ListLength();  // 获取顺序栈长度

    bool IsEmpty();  // 判断顺序栈是否为空
    bool IsFull();   // 判断顺序栈是否已满

private:
    void IncreaseSize();  // 扩容

private:
    T *m_data;      // 存放栈中元素
    int m_maxsize;  // 动态数组最大容量
    int m_top;      // 栈顶指针
};

// 通过构造函数进行初始化
template <typename T>
SeqStack<T>::SeqStack(int length) {
    m_data = new T[length];  // 动态内存分配
    m_maxsize = length;  // 最大容量
    m_top = -1; // 空栈
}

// 通过析构函数释放资源
template <typename T>
SeqStack<T>::~SeqStack() {
    delete[] m_data;
}

template <typename T>
bool SeqStack<T>::Push(const T& e) {
    if (IsFull() == true) {
        IncreaseSize();
    }
    m_top++;
    m_data[m_top] = e;
    return true;
}

// 顺序栈扩容
template <typename T> 
void SeqStack<T>::IncreaseSize() {
    T *p = m_data;
    m_data = new T[m_maxsize + IncSize];

    // 将数据复制到新区域
    for (int i = 0; i <= m_top; i++) {
        m_data[i] = p[i];
    }
    m_maxsize = m_maxsize + IncSize;
    delete[] p;
}

// 出栈
template<typename T>
bool SeqStack<T>::Pop(T& e) {
    if (IsEmpty() == true) {
        std::cout << "Empty Stack" << std::endl;
        return false;
    }
    e = m_data[m_top];
    m_top--;
    return true;
}

template <typename T>
bool SeqStack<T>::GetTop(T &e) {
    if (IsEmpty() == true) {
        std::cout << "Empty Stack" << std::endl;
        return false;
    }
    e = m_data[m_top];
    return true;
}

template <class T>
void SeqStack<T>::DispList() {
    for (int i = m_top; i >= 0; i--) {
        std::cout << m_data[i] << " ";
    }
    std::cout << std::endl;
}

template <class T>
int SeqStack<T>::ListLength() {
    return m_top + 1;
}

template <class T>
bool SeqStack<T>::IsEmpty() {
    if (m_top == -1) {
        return true;
    }
    return false;
}

template <class T>
bool SeqStack<T>::IsFull() {
    if (m_top >= m_maxsize - 1) {
        return true;
    }
    return false;
}

int main(void) {
    SeqStack<int> seqobj(10);
    seqobj.Push(150);
    seqobj.Push(200);
    seqobj.Push(300);
    seqobj.Push(400);
    seqobj.DispList();
    
    int eval = 0;
    seqobj.Pop(eval);
    seqobj.Pop(eval);
    std::cout << "---------------" << std::endl;
    seqobj.DispList();
    
    return 0;
}