#include <iostream>

#define InitSize 10
#define IncSize 5

template <typename T>
class ShareStack {
public:
    ShareStack(int length = InitSize) {
        m_data = new T[length];
        m_maxsize = length;
        m_top1 = -1;
        m_top2 = length;
    }

    ~ShareStack() {
        delete[] m_data;
    }

public:
    bool IsFull() {
        if (m_top1 == m_top2) {
            return true;
        }
        return false;
    }

    bool Push(int stackNum,const T &e) {
        if (IsFull() == true) {
            std::cout << "Full Stack" << std::endl;
            return false;
        }
        if (stackNum == 1) {
            m_top1++;
            m_data[m_top1] = e;
        } else {
            m_top2--;
            m_data[m_top2] = e;
        }
        return true;
    }

    bool Pop(int stackNum,T &e) {
        if (stackNum == 1) {
            if (m_top1 == -1) {
                std::cout << "Share Stack 1 Empty" << std::endl;
                return false;
            }
            e = m_data[m_top1];
            m_top1--;
        } else {
            if (m_top2 == m_maxsize) {
                std::cout << "Share Stack 2 Empty" << std::endl;
                return false;
            }
            e = m_data[m_top2];
            m_top2++;
        }
        return true;
    }

private:
    T *m_data;
    int m_maxsize;
    int m_top1;
    int m_top2;
};

int main(void) {
    ShareStack<int> shareobj(10);
    shareobj.Push(1,150);
    shareobj.Push(2,200);

    int eval2 = 0;
    shareobj.Pop(1,eval2);
    shareobj.Pop(1,eval2);

    return 0;
}